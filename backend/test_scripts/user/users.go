package user

import (
	_ "embed"
	"strings"
)

const testPass = "123"

//go:embed gen.txt
var gen string

type TestUser struct {
	Name  string
	Pass  string
	Email string
}

func TestUsers() []TestUser {
	users := []TestUser{
		{
			Name:  "андрей",
			Email: "importandrew",
		},
		{
			Name:  "五条 悟",
			Email: "stgj",
		},
		{
			Name:  "kokushi",
			Email: "sirota013",
		},
		{
			Name:  "Kokushi",
			Email: "sirota01",
		},
		{
			Name:  "van mao",
			Email: "vanmao",
		},
		{
			Name:  "ivan sidorov",
			Email: "ivan1",
		},
		{
			Name:  "Edwin van der Sar",
			Email: "vratar70",
		},
		{
			Name:  "Edvin wan der Sar",
			Email: "wratar71",
		},
	}

	// fill to 100 users with random users
	for _, s := range strings.Split(gen, "\n") {
		users = append(users, TestUser{Name: s, Email: s})
	}

	for i := range users {
		// all test users' passwords are the same (for front access)
		users[i].Pass = testPass

		// add "@test" to all email fields
		users[i].Email += "@test"
	}

	return users
}
