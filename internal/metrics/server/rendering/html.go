package rendering

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/commons/rendering"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

const (
	MetricsListTemplate = "metrics-list.html"
)

type htmlEngine struct {
	templateParser *rendering.HTMLTemplateParser
}

func NewHTMLEngine(templateParser *rendering.HTMLTemplateParser) *htmlEngine {
	return &htmlEngine{
		templateParser: templateParser,
	}
}

func (e *htmlEngine) RenderList(templateName string, data []domain.Metric) ([]byte, error) {
	return e.templateParser.Parse(templateName, data)
}
