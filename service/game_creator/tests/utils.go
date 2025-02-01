package tests

import "github.com/brianvoe/gofakeit"

func RandomCredentials() (string, string) {
	var username string
	for len(username) < 8 {
		username = gofakeit.Username()
	}
	password := gofakeit.Password(true, true, true, false, false, 8)

	return username, password
}
