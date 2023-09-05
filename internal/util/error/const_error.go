package _error

func BadLogin() error {
	return Err400(map[string][]string{
		"email_or_username": {
			"invalid email or passowrd",
		},
		"password": {
			"invalid email or password",
		},
	})
}

func BadExistField(field string, msg string) error {
	return Err400(map[string][]string{
		field: {
			msg,
		},
	})
}
