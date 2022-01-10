package Server

import (
	"encoding/json"
	"fmt"
	"food-matters/Config"
	"github.com/mattermost/mattermost-server/v5/model"
	"log"
	"net/http"
	"strings"
)

type Command struct {
	UserId    string `schema:"user_id"`
	ChannelId string `schema:"channel_id"`
	TeamId    string `schema:"team_id"`

	ChannelName string `schema:"channel_name"`
	Command     string `schema:"command"`
	ResponseUrl string `schema:"response_url"` //: "https://mm.dev.mrl0co.com/hooks/commands/rpgdf3nco7gp9g198ms6oj3dce"
	TeamDomain  string `schema:"team_domain"`
	Text        string `schema:"text"`
	Token       string `schema:"token"`
	TriggerId   string `schema:"trigger_id"`
	UserName    string `schema:"user_name"`

	SubCommand string   `schema:"-"`
	Args       []string `schema:"-"`
}

type CommandHandler interface {
	GetCommand() string
	Run(command Command) (msg model.CommandResponse, err error)
	HelpText() (msg string)
	HelpDetails(command Command) string
}

var commands = map[string]CommandHandler{}

func AddHandler(command CommandHandler) {
	log.Printf("Registered command: %s", command.GetCommand())
	commands[command.GetCommand()] = command
}

func GetHandler(name string) (CommandHandler, bool) {
	handler, ok := commands[name]
	return handler, ok
}
func GetHandlers() map[string]CommandHandler {
	return commands
}

func HandleCommand(w http.ResponseWriter, req *http.Request) {
	command, err := ParseCommand(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	if command.Token != Config.CommandToken {
		log.Println("unknown source")
		return
	}

	marshal, err := json.MarshalIndent(command, "", "    ")
	log.Printf(string(marshal))

	command.ParseSubCommand()
	commandHandler, ok := commands[command.SubCommand]
	if !ok {
		commandHandler = commands["help"]
	}

	response, err := commandHandler.Run(command)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	fmt.Fprintf(w, response.ToJson())
}

func ParseCommand(req *http.Request) (Command, error) {
	var command Command

	if err := req.ParseForm(); err != nil {
		return command, err
	}

	err := decoder.Decode(&command, req.PostForm)

	return command, err
}

func (command *Command) ParseSubCommand() {
	parts := strings.Fields(command.Text)

	if len(parts) == 0 {
		command.SubCommand = "help"
		return
	}

	command.SubCommand = parts[0]

	if len(parts) > 1 {
		command.Args = parts[1:]
	}
}
