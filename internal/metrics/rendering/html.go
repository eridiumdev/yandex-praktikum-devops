package rendering

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/common/templating"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

const (
	metricsListTemplate = "metrics-list.html"
)

type htmlEngine struct {
	templateParser *templating.HTMLTemplateParser
}

func NewHTMLEngine(templateParser *templating.HTMLTemplateParser) *htmlEngine {
	return &htmlEngine{
		templateParser: templateParser,
	}
}

func (e *htmlEngine) RenderList(list []domain.Metric) ([]byte, error) {
	return e.templateParser.Parse(metricsListTemplate, list)
}
