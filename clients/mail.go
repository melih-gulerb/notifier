package clients

import (
	"context"
	"fmt"

	"github.com/getbrevo/brevo-go/lib"
)

type MailClient struct {
	brevoClient *lib.APIClient
	from        string
}

func NewMailClient(brevoClient *lib.APIClient, from string) *MailClient {
	return &MailClient{
		brevoClient: brevoClient,
		from:        from,
	}
}

func (m *MailClient) SendEmail(toEmail, subject, htmlBody string) {
	sender := lib.SendSmtpEmailSender{
		Name:  "Name",
		Email: m.from,
	}

	to := []lib.SendSmtpEmailTo{
		{
			Email: toEmail,
		},
	}

	sendSmtpEmail := lib.SendSmtpEmail{
		Sender:      &sender,
		To:          to,
		Subject:     subject,
		HtmlContent: htmlBody,
	}

	result, _, err := m.brevoClient.TransactionalEmailsApi.SendTransacEmail(context.Background(), sendSmtpEmail)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Email sent successfully with message ID: %s\n", result.MessageId)
	}
}
