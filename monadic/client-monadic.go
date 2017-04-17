package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rebeccaskinner/agile17-sample/user"
	"github.com/rebeccaskinner/gofpher/either"
)

func main() {
	config, err := getArgs()
	if err != nil {
		fmt.Println(err)
		fmt.Println(showHelp())
		os.Exit(1)
	}

	var (
		getEndpoint  = fmt.Sprintf("%s/oldusers/%s", config.endpoint, config.username)
		postEndpoint = fmt.Sprintf("%s/newusers/%s", config.endpoint, config.username)
		get          = either.WrapEither(http.Get)
	)

	result := get(getEndpoint).
		Next(getResponseBody).
		Next(ioutil.ReadAll).
		Next(user.NewFromJSON).
		Next(user.NewUserFromUser).
		Next(json.Marshal).
		Next(bytes.NewBuffer).
		Next(func(b *bytes.Buffer) (*http.Response, error) {
			return http.Post(postEndpoint, "application/json", b)
		})

	fmt.Println(result)
}

func getResponseBody(r *http.Response) io.Reader {
	return r.Body
}

type config struct {
	endpoint string
	username string
}

func getArgs() (*config, error) {
	args := os.Args
	if len(args) < 3 {
		return nil, errors.New("insufficient number of arguments")
	}
	switch args[1] {
	case "-?", "-h", "--help":
		return nil, errors.New("showing help message")
	}
	return &config{endpoint: args[1], username: args[2]}, nil
}

func showHelp() string {
	return `client: a simple client to get json data about a user
Usage:
client <endpoint> <username>

Client will connect to server at <endpoint> and request data about <username>.

Example:
client http://localhost:8080 user1      # Get information about user1
`
}
