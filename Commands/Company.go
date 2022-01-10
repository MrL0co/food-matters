package Commands

import (
	"food-matters/Config"
	"food-matters/Server"
	"github.com/mattermost/mattermost-server/v5/model"
	"log"
)

type CompanyConfig struct {
}

func (c CompanyConfig) Run(command Server.Command) (msg model.CommandResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (c CompanyConfig) HelpText() (msg string) {
	//TODO implement me
	panic("implement me")
}

func asdasd(command Server.Command) {
	Company := "__NAME_HERE__"
	Menu := "__MENU__"

	dialog := model.OpenDialogRequest{
		TriggerId: command.TriggerId,
		URL:       Config.GetServerAddress() + "/hello",
		Dialog: model.Dialog{
			CallbackId:       "some random id",
			Title:            "Order your food bro",
			IntroductionText: "Order from: " + Company + " menu: " + Menu,
			IconURL:          "",
			Elements: []model.DialogElement{
				{
					DisplayName: "food",
					Name:        "food",
					Type:        "text",
					SubType:     "text",
					Placeholder: "food name",
					HelpText:    "",
					Optional:    false,
					MinLength:   3,
					MaxLength:   0,
				},
			},
			SubmitLabel:    "Order",
			NotifyOnCancel: false,
			State:          "state test",
		},
	}

	if success, resp := Config.Client.OpenInteractiveDialog(dialog); !success {
		log.Println(resp)
	} else {
		log.Println("send dialog")
	}

}
