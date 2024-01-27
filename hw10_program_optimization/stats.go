package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrCannotUnmarshalJSON = errors.New("cannot unmarshal json")

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	br := bufio.NewReader(r)
	var user User
	i := 0
	for {
		line, _, readErr := br.ReadLine()

		if readErr == io.EOF {
			break
		}

		if err = user.UnmarshalJSON(line); err != nil {
			return result, ErrCannotUnmarshalJSON
		}

		result[i] = user
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := strings.Contains(user.Email, "."+domain)

		if matched {
			subDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[subDomain]++
		}
	}
	return result, nil
}
