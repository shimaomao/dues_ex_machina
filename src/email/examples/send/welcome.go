package send

import (
	"github.com/matcornic/hermes/v2"
)

type welcome struct {
}

func (w *welcome) Name() string {
	return "welcome"
}

func (w *welcome) Email(urlToken string, endPoint string, targetEmail string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: "Jon Snow",
			Intros: []string{
				"Welcome to Hermes! We're very excited to have you on board.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Firstname", Value: "Jon"},
				{Key: "Lastname", Value: "Snow"},
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Hermes, please click here:",
					Button: hermes.Button{
						Text: "Confirm your account",
						Link: endPoint + "/api/user/confirm?token=" + urlToken + "&email=" + targetEmail,
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
}
