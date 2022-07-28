package templating

import (
	"bytes"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
)

type HTMLTemplateParser struct {
	tpl *template.Template
}

func NewHTMLTemplateParser(templatesDir string) *HTMLTemplateParser {
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		logger.Fatalf("[html template parser] directory '%s' does not exist", templatesDir)
	}

	var tpl *template.Template
	err := filepath.WalkDir(templatesDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			tpl, err = template.ParseFiles(path)
		}
		return err
	})
	if err != nil {
		logger.Fatalf("[html template parser] error when parsing templates: %s", err.Error())
	}

	return &HTMLTemplateParser{
		tpl: tpl,
	}
}

func (r *HTMLTemplateParser) Parse(templateName string, data any) ([]byte, error) {
	var buffer bytes.Buffer

	err := r.tpl.ExecuteTemplate(&buffer, templateName, data)
	if err != nil {
		return nil, errors.Wrapf(err, "[html template parser]")
	}
	return buffer.Bytes(), nil
}
