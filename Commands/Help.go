package Commands

import (
	"fmt"
	"food-matters/Markdown"
	"food-matters/Server"
	"food-matters/Translations"
	"github.com/mattermost/mattermost-server/v5/model"
)

type Help struct {
	Name string
}

func NewHelp() *Help {
	return &Help{Name: Translations.SLocalize(Translations.HelpName)}
}

func (h Help) GetCommand() string {
	return h.Name
}

func (h Help) Run(command Server.Command) (model.CommandResponse, error) {
	var message = ""
	var unknownCommand = ""

	if command.SubCommand != h.GetCommand() {
		unknownCommand = command.SubCommand
	} else if len(command.Args) > 0 {
		if handler, ok := Server.GetHandler(command.Args[0]); ok {
			message = handler.HelpDetails(command)
		} else {
			unknownCommand = command.Args[0]
		}
	}

	if message == "" {
		message = h.HelpDetails(command)
	}

	if unknownCommand != "" {
		message = fmt.Sprintf(
			"%s\n%s",
			Translations.Localize(Translations.UnknownCommand, map[string]interface{}{"Command": unknownCommand}),
			message,
		)
	}

	return model.CommandResponse{
		ResponseType: "ephemeral",
		Text:         message,
	}, nil
}

func (h Help) HelpText() (msg string) {
	return Translations.SLocalize(Translations.HelpHelpText)
}

func (h Help) HelpDetails(command Server.Command) (msg string) {
	t := Markdown.NewTable("subcommand", "description")

	for name, command := range Server.GetHandlers() {
		t.AddRow(name, command.HelpText())
	}

	return "Usage: `" + command.Command + " <subcommand> [args]`\n" +
		"Food Ordering system for mattermost\n" +
		"Type `" + command.Command + " help <subcommand>` for help with a specific subcommand.\n" +
		"available subcommands:\n\n" + t.ToString()
}
