package helm

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/storage/driver"
	"sigs.k8s.io/yaml"
)

var (
	settings = cli.New()
)

type repoAddOptions struct {
	name     string
	url      string
	username string
	password string
	noUpdate bool

	certFile string
	keyFile  string
	caFile   string

	repoFile  string
	repoCache string
}

func debug(format string, v ...interface{}) {
	if settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		logger.Debug(fmt.Sprintf(format, v...))
	}
}

// PullChart pulls a chart from helm registry
func PullChart(chartName string, chartVersion string) error {
	logger.Debug("New pull client")
	actionConfig := new(action.Configuration)
	actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug)

	client := action.NewPull()
	client.Settings = settings
	client.Version = chartVersion

	logger.Debugf("Pulling %s:%s...", chartName, chartVersion)
	res, err := client.Run(chartName)
	if err != nil {
		return err
	}
	logger.Infof("Pulled %s:%s %s", chartName, chartVersion, res)

	return nil
}

// PullChartUntarToDir pulls a helm chart, untar and output to a directory
func PullChartUntarToDir(chartName string, chartVersion string, dirName string) error {
	logger.Debug("New pull client")
	actionConfig := new(action.Configuration)
	actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug)

	client := action.NewPull()
	client.Settings = settings
	client.Version = chartVersion
	client.Untar = true
	client.UntarDir = dirName

	logger.Debugf("Pulling %s:%s...", chartName, chartVersion)
	res, err := client.Run(chartName)
	if err != nil {
		return err
	}
	logger.Infof("Pulled %s:%s %s", chartName, chartVersion, res)

	return nil
}

// AddRepository adds a repository to helm
func AddRepository(repoName string, repoURL string) error {
	logger.Debug("New repository client")
	actionConfig := new(action.Configuration)
	actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug)

	opts := &repoAddOptions{}

	opts.name = repoName
	opts.url = repoURL
	opts.repoFile = settings.RegistryConfig
	opts.repoCache = settings.RepositoryCache

	return opts.run()
}

func (opts *repoAddOptions) run() error {
	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(opts.repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(opts.repoFile, filepath.Ext(opts.repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(opts.repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if opts.noUpdate && f.Has(opts.name) {
		return errors.Errorf("repository name (%s) already exists, please specify a different name", opts.name)
	}

	c := repo.Entry{
		Name:     opts.name,
		URL:      opts.url,
		Username: opts.username,
		Password: opts.password,
		CertFile: opts.certFile,
		KeyFile:  opts.keyFile,
		CAFile:   opts.caFile,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", opts.url)
	}

	f.Update(&c)

	if err := f.WriteFile(opts.repoFile, 0644); err != nil {
		return err
	}
	logger.Infof("%q has been added to your repositories\n", opts.name)

	return nil
}

// UpgradeInstallChart upgrades a existing release or creates it
func UpgradeInstallChart(releaseName string, chartPath string, valueOpts *values.Options) error {
	logger.Debug("New upgrade-install client")
	actionConfig := new(action.Configuration)
	actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug)

	client := action.NewUpgrade(actionConfig)

	p := getter.All(settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return err
	}

	chartPath, err = client.ChartPathOptions.LocateChart(chartPath, settings)
	if err != nil {
		return err
	}

	// If a release does not exist, install it. If another error occurs during
	// the check, ignore the error and continue with the upgrade.
	histClient := action.NewHistory(actionConfig)
	histClient.Max = 1
	if _, err := histClient.Run(releaseName); err == driver.ErrReleaseNotFound {
		logger.Infof("Release %q does not exist. Installing it now.\n", releaseName)
		instClient := action.NewInstall(actionConfig)
		instClient.ReleaseName = releaseName
		instClient.ChartPathOptions = client.ChartPathOptions
		instClient.DryRun = client.DryRun
		instClient.DisableHooks = client.DisableHooks
		instClient.Timeout = client.Timeout
		instClient.Wait = client.Wait
		instClient.Devel = client.Devel
		instClient.Namespace = client.Namespace
		instClient.Atomic = client.Atomic

		err := runInstall(instClient, releaseName, chartPath, valueOpts)
		if err != nil {
			return err
		}
	}

	// Check chart dependencies to make sure all are present in /charts
	ch, err := loader.Load(chartPath)
	if err != nil {
		return err
	}
	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return err
		}
	}

	if ch.Metadata.Deprecated {
		logger.Warning("WARNING: This chart is deprecated")
	}

	logger.Infof("Installing %s in %s", releaseName, client.Namespace)
	_, err = client.Run(releaseName, ch, vals)
	if err != nil {
		return errors.Wrap(err, "UPGRADE FAILED")
	}

	return nil
}

// runInstall installs a helm chart in current context
func runInstall(client *action.Install, releaseName string, chartPath string, valueOpts *values.Options) error {
	logger.Debug("New install client")

	cp, err := client.ChartPathOptions.LocateChart(chartPath, settings)
	if err != nil {
		return err
	}

	p := getter.All(settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return err
	}

	if chartRequested.Metadata.Deprecated {
		logger.Warning("This chart is deprecated")
	}

	client.Namespace = settings.Namespace()
	_, err = client.Run(chartRequested, vals)
	if err != nil {
		return err
	}

	return nil
}

// isChartInstallable validates if a chart can be installed
//
// Application chart type is only installable
func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
