package Translations

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	ServerStarted = &i18n.Message{
		ID:    "ServerStarted",
		Other: "_{{.Name}} has **started** running_",
	}
	HelpName = &i18n.Message{
		ID:    "HelpName",
		Other: "help",
	}
	HelpHelpText = &i18n.Message{
		ID:    "HelpName",
		Other: "Show help message",
	}

	UnknownCommand = &i18n.Message{
		ID:    "UnknownCommand",
		Other: "**Unknown command: '{{.Command}}'**",
	}
)
