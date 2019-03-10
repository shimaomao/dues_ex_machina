package utils

import (
	"github.com/domodwyer/mailyak"
	"net/smtp"
)

// StartTLS Email Example

func Email() {
	mail := mailyak.New("smtp-mail.outlook.com:587", smtp.PlainAuth("", "kecy6tgy@nottingham.edu.my", "Ilovecoffee96@", "smtp-mail.outlook.com"))

	mail.To("alex.tay@jonvi.com")
	mail.From("kecy6tgy@nottingham.edu.my")
	mail.FromName("Prince Fournineteen")

	mail.Subject("Business proposition")

	// mail.HTML() and mail.Plain() implement io.Writer, so you can do handy things like
	// parse a template directly into the email body
	// if err := t.ExecuteTemplate(mail.HTML(), "htmlEmail", data); err != nil {
	// 	panic(" 💣 ")
	// }

	// Or set the body using a string setter
	mail.Plain().Set("Get a real email client")

	// And you're done!
	if err := mail.Send(); err != nil {
		panic(" 💣 ")
	}

}
