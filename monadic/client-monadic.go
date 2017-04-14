package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rebeccaskinner/agile17-sample/user"
	"github.com/rebeccaskinner/gofpher/either"
	"github.com/rebeccaskinner/gofpher/monad"
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
		read         = either.WrapEither(ioutil.ReadAll)
		fromjson     = either.WrapEither(user.NewFromJSON)
		mkUser       = either.WrapEither(user.NewUserFromUser)
		toJSON       = either.WrapEither(json.Marshal)
	)

	usr := get(getEndpoint).
		AndThen(read).
		AndThen(fromjson).
		AndThen(mkUser).
		AndThen(toJSON)

	buffered := monad.FMap(bytes.NewBuffer, usr).(either.EitherM)

	if buffered.IsLeft() {
		fmt.Println(buffered.FromLeft().(error))
		os.Exit(1)
	}

	buf := buffered.FromRight().(*bytes.Buffer)

	response, err := http.Post(postEndpoint, "application/json", buf)

	if err != nil {
		fmt.Println("failed to post message: ", err)
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Printf("failed to post message: server returned: " + response.Status)
		os.Exit(1)
	}

	fmt.Printf("Success")

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
