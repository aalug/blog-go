package mail

import (
	"github.com/aalug/blog-go/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithMailHog(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewHogSender(config.EmailSenderAddress)

	subject := "A test email"
	content := `
	<h1>Hello</h1>
	<p>This is a test message</p>
	`
	to := []string{"test@example.com"}
	attachFiles := []AttachFile{
		{
			Name: "readme file",
			Path: "../README.md",
		},
	}

	err = sender.SendEmail(Data{
		To:      to,
		Subject: subject,
		Content: content,
		Files:   attachFiles,
	})
	require.NoError(t, err)
}
