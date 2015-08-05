package account

// parse aws credentials formatted file and create a list of accounts
// to iterate over

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Account struct {
	Name   string
	Key    string
	Secret string
}

func (a *Account) Done() bool {
	return (a.Name != "" && a.Key != "" && a.Secret != "")
}

func (a *Account) AddName(s string) {
	a.Name = strings.Replace(s, "]", "", -1)
	a.Name = strings.Replace(a.Name, "[", "", -1)
}

func (a *Account) AddKey(s string) {
	parts := strings.Split(s, "=")
	for n, p := range parts {
		parts[n] = strings.TrimSpace(p)
	}

	a.Key = parts[1]
}

func (a *Account) AddSecret(s string) {
	parts := strings.Split(s, "=")
	for n, p := range parts {
		parts[n] = strings.TrimSpace(p)
	}

	a.Secret = parts[1]
}

func ParseAccounts(f *string) (accts []*Account) {
	in, err := ioutil.ReadFile(*f)
	if err != nil {
		panic(err)
	}

	file := strings.Split(string(in), "\n")

	if len(file) == 0 {
		panic(fmt.Errorf("expected data, got none"))
	}

	acct := new(Account)
	for _, line := range file {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") {
			acct.AddName(line)
		}

		if strings.HasPrefix(line, "aws_access_key_id") {
			acct.AddKey(line)
		}

		if strings.HasPrefix(line, "aws_secret_access_key") {
			acct.AddSecret(line)
		}

		if acct.Done() {
			accts = append(accts, acct)
			acct = new(Account)
		}
	}

	return accts
}
