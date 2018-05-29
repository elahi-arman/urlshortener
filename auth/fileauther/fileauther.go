package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/elahi-arman/urlshortener/auth"
)

type fileauth struct {
	fh *os.File
}

//Terminate closes the file handler
func (f fileauth) Terminate() {
	f.fh.Close()
}

// NewFileAuth creates a new file auth
func NewFileAuth(fpath string) auth.Auther {
	file, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}

	return &fileauth{
		fh: file,
	}

}

//IsValid determines if the user is valid
func (f fileauth) IsValid(user string, pass string) bool {
	scanner := bufio.NewScanner(f.fh)
	for scanner.Scan() {
		line := scanner.Text()
		authStr := fmt.Sprintf("%s:%s", user, pass)

		if strings.Contains(line, authStr) {
			return true
		}
	}

	return false
}

//IsUserAuthorizedForScope verifies if a user is allowed to r/w the scope
func (f fileauth) IsUserAuthorizedForScope(user string, scope string) bool {
	scanner := bufio.NewScanner(f.fh)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, user) {
			userline := strings.SplitN(line, ":", 3)
			scopes := userline[2]
			for _, s := range strings.Split(scopes, ",") {
				if strings.Compare(s, scope) == 0 {
					return true
				}
			}

		}
	}

	return false
}
