package services

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/mailer"
)

type IMailerService interface {
	SendMailForgotPassword(user *models.User) error
}

var config = mailer.GomailSenderConfig{
	Host:     utils.GetEnv("MAIL_HOST", "smtp.gmail.com"),
	Port:     utils.GetEnvAsInt("MAIL_PORT", 587),
	Username: utils.GetEnv("MAIL_USERNAME", "vfa.khuongdv@gmail.com"),
	Password: utils.GetEnv("MAIL_PASSWORD", "hupr ojqr nwkq tuzo"),
	From:     utils.GetEnv("MAIL_FORM", "vfa.khuongdv@gmail.com"),
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
	mailer := mailer.NewGomailSender(mailer.GomailSenderConfig{
		From:     config.From,
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	})

	// Parse the email template file
	tmpl, err := template.ParseFiles("pkg/mailer/templates/forgot_template.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Construct reset password URL by combining frontend URL with user's reset token
	url := utils.GetEnv("FRONTEND_URL", "") + "/reset-password?token=" + *user.Token

	// Prepare template data with user's name and reset URL
	data := map[string]interface{}{
		"Name": user.Name,
		"URL":  url,
	}
	// Create buffer to store rendered HTML
	var htmlBody bytes.Buffer
	// Execute template with data and write to buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Send password reset email to user
	if err := mailer.Send([]string{user.Email}, "Reset your password", "", htmlBody.String()); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil

}
