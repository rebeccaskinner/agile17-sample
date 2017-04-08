package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rebeccaskinner/agile17-sample/user"
)

func main() {
	config, err := getArgs()
	if err != nil {
		fmt.Println(err)
		fmt.Println(showHelp())
		os.Exit(1)
	}

	response, err := http.Get(config.endpoint + "/fetch/" + config.username)
	if err != nil {
		fmt.Println("unable to fetch HTTP data", err)
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println("server returned " + response.Status)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("failed to read body: ", err)
		os.Exit(1)
	}

	user, err := user.NewFromJSON(body)
	if err != nil {
		fmt.Println("failed to deserialize json: ", err)
		os.Exit(1)
	}

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
