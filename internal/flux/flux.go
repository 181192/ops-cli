package flux

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/shurcooL/httpfs/vfsutil"
)

//go:generate go run generate.go

// TemplateParameters represent template parameters
type TemplateParameters struct {
	GitURL             string
	GitBranch          string
	GitPaths           []string
	GitLabel           string
	GitUser            string
	GitEmail           string
	GitReadOnly        bool
	ManifestGeneration bool
	GarbageCollection  bool
	AcrRegistry        bool
	HelmVersions       []string
}

// FillInTemplates fils in custom helm values for flux and helm-operator
func FillInTemplates(params *TemplateParameters) (map[string][]byte, error) {
	result := map[string][]byte{}
	err := vfsutil.WalkFiles(templates, string(os.PathSeparator), func(path string, info os.FileInfo, rs io.ReadSeeker, err error) error {
		if err != nil {
			return fmt.Errorf("cannot walk embedded files: %s", err)
		}

		if info.IsDir() {
			return nil
		}

		manifestTemplateBytes, err := ioutil.ReadAll(rs)
		if err != nil {
			return fmt.Errorf("cannot read embedded file %q: %s", info.Name(), err)
		}

		manifestTemplate, err := template.New(info.Name()).
			Funcs(template.FuncMap{"StringsJoin": strings.Join}).
			Parse(string(manifestTemplateBytes))
		if err != nil {
			return fmt.Errorf("cannot parse embedded file %q: %s", info.Name(), err)
		}

		out := bytes.NewBuffer(nil)
		if err := manifestTemplate.Execute(out, params); err != nil {
			return fmt.Errorf("cannot execute template for embedded file %q: %s", info.Name(), err)
		}

		if out.Len() > 0 {
			result[strings.TrimSuffix(info.Name(), ".tmpl")] = out.Bytes()
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("internal error filling embedded installation templates: %s", err)
	}

	return result, nil
}
