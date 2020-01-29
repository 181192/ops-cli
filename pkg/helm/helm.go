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
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
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
