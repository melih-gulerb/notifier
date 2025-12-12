package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

func RenderTemplate(templateCode string, data TemplateData) (string, string, error) {
	emailTemplate, exists := templateMap[templateCode]
	if !exists {
		return "", "", fmt.Errorf("template not found: %s", templateCode)
	}

	subjectTmpl, err := template.New("subject").Parse(emailTemplate.Subject)
	if err != nil {
		return "", "", err
	}

	var subjectBuf bytes.Buffer
	if err := subjectTmpl.Execute(&subjectBuf, data); err != nil {
		return "", "", err
	}

	bodyTmpl, err := template.New("body").Parse(emailTemplate.Body)
	if err != nil {
		return "", "", err
	}

	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, data); err != nil {
		return "", "", err
	}

	return subjectBuf.String(), bodyBuf.String(), nil
}
