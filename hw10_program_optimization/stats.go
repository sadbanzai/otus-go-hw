package hw10programoptimization

import (
	"bufio"
	"errors"
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
	br := bufio.NewReader(r)
	var user User
	result := make(DomainStat)
	for {
		line, _, readErr := br.ReadLine()

		if readErr == io.EOF {
			break
		}

		if err := user.UnmarshalJSON(line); err != nil {
			return result, ErrCannotUnmarshalJSON
		}

		matched := strings.HasSuffix(user.Email, "."+domain)

		if matched {
			subDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[subDomain]++
		}
	}
	return result, nil
}
