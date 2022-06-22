package rendering

type Engine interface {
	Render(templatePath string, data any) ([]byte, error)
}
