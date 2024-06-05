package validation

import "fmt"

func Verb(verb string) error {
	if verb == "" {
		return fmt.Errorf("empty verb")
	}
	return nil
}

func Account(account string) error {
	if account == "" {
		return fmt.Errorf("empty account")
	}
	return nil
}
