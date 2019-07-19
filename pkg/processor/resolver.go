package processor

// TemplateResolver interface defines how to find template
type TemplateResolver interface {

	Locate(template string, locale string) *Template

}

type defaultResolver struct {

}

// Template struct
type Template struct {
	Raw string
	Values map[string]interface{}
	Compiled bool
}

func (a *defaultResolver) Locate(template string, locale string) *Template {
	return &Template{
		Raw: "Some text",
		Compiled: false,
	}
}