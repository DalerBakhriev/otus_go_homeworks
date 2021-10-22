package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

type users []User

func getUsers(r io.Reader) (users, error) {
	scanner := bufio.NewScanner(r)

	i := 0
	user := User{}
	result := make(users, 0, 100_000)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return result, err
		}
		result = append(result, user)
		i++
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainCheck := strings.Join([]string{".", domain}, "")
	for _, user := range u {
		if strings.Contains(user.Email, domainCheck) {
			emailDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[emailDomain]
			num++
			result[emailDomain] = num
		}
	}
	return result, nil
}
