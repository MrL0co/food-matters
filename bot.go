// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package main

import (
	"food-matters/Commands"
	"food-matters/Config"
	"food-matters/Server"
	"food-matters/Translations"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

//	_Food Matters has **stopped** running_
func main() {
	log.Println(Config.Name)

	Server.AddHandler(Commands.NewHelp())

	SetupGracefulShutdown()

	Config.Client = model.NewAPIv4Client(Config.ServerUrl)

	MakeSureServerIsRunning()
	LoginAsTheBotUser()
	UpdateTheBotUserIfNeeded()
	FindBotTeam()
	CreateBotDebuggingChannelIfNeeded()

	Server.SendMsgToDebuggingChannel(Translations.Localize(Translations.ServerStarted, map[string]interface{}{"Name": Config.Name}), "")

	webSocketClient, appErr := model.NewWebSocketClient4(Config.GetWebsocketServer(), Config.Client.AuthToken)
	if appErr != nil {
		log.Println("We failed to connect to the web socket")
		Server.PrintError(appErr)
	}

	webSocketClient.Listen()

	go func() {
		for resp := range webSocketClient.EventChannel {
			HandleWebSocketResponse(resp)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/command", Server.HandleCommand)

	log.Println(Config.GetOutboundIP().String())
	err := http.ListenAndServe(Config.GetPort(), mux)
	log.Fatal(err)
}

func hello(w http.ResponseWriter, req *http.Request) {
	dialogResponse, err := Server.ParseDialogResponse(w, req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(dialogResponse)
}

func HandleWebSocketResponse(event *model.WebSocketEvent) {
	if Config.Debug {
		HandleMsgFromDebuggingChannel(event)
	}

	HandleMsg(event)
}

func HandleMsg(event *model.WebSocketEvent) {
	// If this isn't the debugging channel then lets ignore it
	if event.GetBroadcast().ChannelId != Config.DebuggingChannel.Id {
		return
	}

	// Lets only responded to messaged posted events
	if event.EventType() != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	post := model.PostFromJson(strings.NewReader(event.GetData()["post"].(string)))
	if post != nil {
		// ignore my events
		if post.UserId == Config.BotUser.Id {
			return
		}
		// ignore System events
		if post.IsSystemMessage() {
			return
		}

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)show_form(?:$|\W)`, post.Message); matched {

		}
	}
}

func HandleMsgFromDebuggingChannel(event *model.WebSocketEvent) {
	// If this isn't the debugging channel then lets ingore it
	if event.GetBroadcast().ChannelId != Config.DebuggingChannel.Id {
		return
	}

	// Lets only responded to messaged posted events
	if event.EventType() != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	log.Println("responding to debugging channel msg")

	post := model.PostFromJson(strings.NewReader(event.GetData()["post"].(string)))
	if post != nil {

		// ignore my events
		if post.UserId == Config.BotUser.Id {
			return
		}

		// ignore System events
		if post.IsSystemMessage() {
			return
		}

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)(?:alive|up|running|hello)(?:$|\W)`, post.Message); matched {
			Server.SendMsgToDebuggingChannel("Yes I'm running", post.Id)
			return
		}
	}

	Server.SendMsgToDebuggingChannel("I did not understand you!", post.Id)
}

func SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if Config.WebSocketClient != nil {
				Config.WebSocketClient.Close()
			}

			Server.SendMsgToDebuggingChannel("_"+Config.Name+" has **stopped** running_", "")
			os.Exit(0)
		}
	}()
}
