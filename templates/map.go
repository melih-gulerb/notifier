package templates

import "notifier/templates/html"

type TemplateData map[string]interface{}

type EmailTemplate struct {
	Subject string
	Body    string
}

var templateMap = map[string]EmailTemplate{
	"TEST_TEMPLATE": {
		Subject: "Test Subject",
		Body:    html.TestTemplate,
	},
}
