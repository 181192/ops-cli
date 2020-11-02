package fileprocessor

// NoOpTemplateProcessor is a FileProcessor that does no operations on templated files and ignore them from output
type NoOpTemplateProcessor struct{}

// ProcessFile takes a template file and executes the template applying the TemplateParameters
func (p *NoOpTemplateProcessor) ProcessFile(file File) (File, error) {
	return file, nil
}
