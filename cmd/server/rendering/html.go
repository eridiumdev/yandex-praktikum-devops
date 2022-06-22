package rendering

import (
	"bytes"
	"github.com/pkg/errors"
	"html/template"
)

type HTMLEngine struct {
	TemplatesDir string
}

var (
	ErrCannotParseTemplate   = errors.New("cannot load or parse template")
	ErrCannotExecuteTemplate = errors.New("cannot execute template")
)

func NewHTMLEngine(tplDir string) *HTMLEngine {
	return &HTMLEngine{
		TemplatesDir: tplDir,
	}
}

func (e *HTMLEngine) Render(templatePath string, data any) ([]byte, error) {
	t := template.New("index.html")
	t, err := t.ParseFiles(e.TemplatesDir + "/" + templatePath)
	if err != nil {
		return nil, errors.Wrapf(ErrCannotParseTemplate, "[html engine]")
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)
	if err != nil {
		return nil, errors.Wrapf(ErrCannotExecuteTemplate, "[html engine]")
	}
	return buffer.Bytes(), nil
}
