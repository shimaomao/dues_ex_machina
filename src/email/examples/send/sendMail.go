package send

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"

	"github.com/go-gomail/gomail"
	"github.com/matcornic/hermes/v2"
)

type example interface {
	Email() hermes.Email
	Name() string
}

// func GenerateSmtpConfig() (smtpConfig smtpAuthentication, options sendOptions) {

// 	return smtpConfig, options
// }

func SendEmailResetPassword() {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "Hermes",
			Link: "https://example-hermes.com/",
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}
	template := new(reset)

	theme := new(hermes.Default)

	h.Theme = theme

	generateEmails(h, template.Email(), template.Name())
	// smtpConfig, options = GenerateSmtpConfig()
	port := 465
	password := `ezbot5512`
	SMTPUser := "support@ezbot.my"

	smtpConfig := smtpAuthentication{
		Server:         "sg2plcpnl0096.prod.sin2.secureserver.net",
		Port:           port,
		SenderEmail:    "support@ezbot.my",
		SenderIdentity: "Alex Tay",
		SMTPPassword:   password,
		SMTPUser:       SMTPUser,
	}
	options := sendOptions{
		To: "alex.tay@jonvi.com",
	}
	options.Subject = "Testing Hermes - Theme " + h.Theme.Name() + " - Example " + template.Name()
	fmt.Printf("Sending email '%s'...\n", options.Subject)
	htmlBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/%v.%v.html", h.Theme.Name(), h.Theme.Name(), template.Name()))
	if err != nil {
		fmt.Println(err)
	}
	txtBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/%v.%v.txt", h.Theme.Name(), h.Theme.Name(), template.Name()))
	if err != nil {
		fmt.Println(err)
	}
	err = send(smtpConfig, options, string(htmlBytes), string(txtBytes))
	if err != nil {
		fmt.Println(err)
	}
}

func SendEmailVerification() {

	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "Hermes",
			Link: "https://example-hermes.com/",
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}
	template := new(welcome)

	// examples := []example{
	// 	new(welcome),
	// 	// new(reset),
	// 	// new(receipt),
	// 	// new(maintenance),
	// }
	theme := new(hermes.Default)

	// themes := []hermes.Theme{
	// 	new(hermes.Default),
	// }
	h.Theme = theme

	generateEmails(h, template.Email(), template.Name())
	// smtpConfig, options = GenerateSmtpConfig()
	port := 465
	password := `ezbot5512`
	SMTPUser := "support@ezbot.my"

	smtpConfig := smtpAuthentication{
		Server:         "sg2plcpnl0096.prod.sin2.secureserver.net",
		Port:           port,
		SenderEmail:    "support@ezbot.my",
		SenderIdentity: "Alex Tay",
		SMTPPassword:   password,
		SMTPUser:       SMTPUser,
	}
	options := sendOptions{
		To: "alex.tay@jonvi.com",
	}
	options.Subject = "Testing Hermes - Theme " + h.Theme.Name() + " - Example " + template.Name()
	fmt.Printf("Sending email '%s'...\n", options.Subject)
	htmlBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/%v.%v.html", h.Theme.Name(), h.Theme.Name(), template.Name()))
	if err != nil {
		fmt.Println(err)
	}
	txtBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/%v.%v.txt", h.Theme.Name(), h.Theme.Name(), template.Name()))
	if err != nil {
		fmt.Println(err)
	}
	err = send(smtpConfig, options, string(htmlBytes), string(txtBytes))
	if err != nil {
		fmt.Println(err)
	}

}

func generateEmails(h hermes.Hermes, email hermes.Email, example string) {
	// Generate the HTML template and save it
	res, err := h.GenerateHTML(email)
	if err != nil {
		fmt.Println(err)

	}
	err = os.MkdirAll(h.Theme.Name(), 0744)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("%v/%v.%v.html", h.Theme.Name(), h.Theme.Name(), example), []byte(res), 0644)
	if err != nil {
		fmt.Println(err)
	}

	// Generate the plaintext template and save it
	res, err = h.GeneratePlainText(email)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("%v/%v.%v.txt", h.Theme.Name(), h.Theme.Name(), example), []byte(res), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

type smtpAuthentication struct {
	Server         string
	Port           int
	SenderEmail    string
	SenderIdentity string
	SMTPUser       string
	SMTPPassword   string
}

// sendOptions are options for sending an email
type sendOptions struct {
	To      string
	Subject string
}

// send sends the email
func send(smtpConfig smtpAuthentication, options sendOptions, htmlBody string, txtBody string) error {

	if smtpConfig.Server == "" {
		return errors.New("SMTP server config is empty")
	}
	if smtpConfig.Port == 0 {
		return errors.New("SMTP port config is empty")
	}

	if smtpConfig.SMTPUser == "" {
		return errors.New("SMTP user is empty")
	}

	if smtpConfig.SenderIdentity == "" {
		return errors.New("SMTP sender identity is empty")
	}

	if smtpConfig.SenderEmail == "" {
		return errors.New("SMTP sender email is empty")
	}

	if options.To == "" {
		return errors.New("no receiver emails configured")
	}

	from := mail.Address{
		Name:    smtpConfig.SenderIdentity,
		Address: smtpConfig.SenderEmail,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("To", options.To)
	m.SetHeader("Subject", options.Subject)

	m.SetBody("text/plain", txtBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(smtpConfig.Server, smtpConfig.Port, smtpConfig.SMTPUser, smtpConfig.SMTPPassword)

	return d.DialAndSend(m)
}
