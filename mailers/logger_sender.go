package mailers

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo/mail"
)

// LoggerSender is a mail.Sender that simply logs the email to the console.
type LoggerSender struct{}

// Send logs the email
func (l LoggerSender) Send(m mail.Message) error {
	fmt.Printf("\n--- MOCK EMAIL ---\n")
	fmt.Printf("From: %s\n", m.From)
	fmt.Printf("To: %s\n", strings.Join(m.To, ", "))
	fmt.Printf("Subject: %s\n", m.Subject)

	if len(m.Bodies) > 0 {
		fmt.Printf("Body:\n%s\n", m.Bodies[0].Content)
	}

	fmt.Printf("------------------\n\n")
	return nil
}
