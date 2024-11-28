package services

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/mailer"
)

var config = mailer.SMTPConfig{
	Host:     utils.GetEnv("MAIL_HOST", "smtp"),
	Port:     utils.GetEnvAsInt("MAIL_PORT", 587),
	Username: utils.GetEnv("MAIL_USERNAME", ""),
	Password: utils.GetEnv("MAIL_PASSWORD", ""),
	From:     utils.GetEnv("MAIL_FORM", "noreply@vfa-spl.org"),
}

// SendMailForgotPassword sends a password reset email to the user
// Parameters:
//   - user: Pointer to models.User containing user information including email and reset token
//
// Returns:
//   - error: Returns nil on success, error on failure
//
// The function:
//  1. Creates SMTP config from environment variables
//  2. Initializes mail sender
//  3. Parses email template
//  4. Executes template with user data
//  5. Sends password reset email to user
func SendMailForgotPassword(user *models.User) error {
	mailer := mailer.NewGomailSender(mailer.SMTPConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		From:     config.From,
	})

	tmpl, err := template.ParseFiles("forgot_template.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Execute the template with the provided data
	url := utils.GetEnv("FRONTEND_URL", "") + "/reset-password/" + *user.Token

	data := map[string]interface{}{
		"Name": user.Name,
		"URL":  url,
	}
	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	if err := mailer.Send([]string{user.Email}, "Reset your password", htmlBody.String()); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil

}
