package fileprocessor

import (
	"bytes"
	"strings"
	"text/template"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
)

const (
	templateExtension = ".tmpl"
)

// TemplateParameters represents the API variables that can be used to template a profile. This set of variables will
// be applied to the go template files found. Templates filenames must end in .templ
type TemplateParameters struct {
	ClusterName string
	Location    string
}

// NewTemplateParameters creates a set of variables for templating given a ClusterConfig object
func NewTemplateParameters(clusterConfig *api.AKSClusterConfig) TemplateParameters {
	return TemplateParameters{
		ClusterName: clusterConfig.ObjectMeta.Name,
		Location:    clusterConfig.Spec.Location,
	}
}

// GoTemplateProcessor is a FileProcessor that executes Go Templates
type GoTemplateProcessor struct {
	Params TemplateParameters
}

// ProcessFile takes a template file and executes the template applying the TemplateParameters
func (p *GoTemplateProcessor) ProcessFile(file File) (File, error) {
	if !isGoTemplate(file.Path) {
		logger.Debug("leaving non-template file unmodified %q", file.Path)
		return file, nil
	}

	parsedTemplate, err := template.New(file.Path).Parse(string(file.Data))
	if err != nil {
		return File{}, errors.Wrapf(err, "cannot parse manifest template file %q", file.Path)
	}

	out := bytes.NewBuffer(nil)
	if err = parsedTemplate.Execute(out, p.Params); err != nil {
		return File{}, errors.Wrapf(err, "cannot execute template for file %q", file.Path)
	}

	newFileName := strings.TrimSuffix(file.Path, templateExtension)
	return File{
		Data: out.Bytes(),
		Path: newFileName,
	}, nil
}

func isGoTemplate(fileName string) bool {
	return strings.HasSuffix(fileName, templateExtension)
}
