package main

import (
	"food-matters/Config"
	"food-matters/Server"
	"github.com/mattermost/mattermost-server/v5/model"
	"log"
	"os"
)

func MakeSureServerIsRunning() {
	if props, resp := Config.Client.GetOldClientConfig(""); resp.Error != nil {
		log.Println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		Server.PrintError(resp.Error)
		os.Exit(1)
	} else {
		log.Println("Server detected and is running version " + props["Version"])
	}
}

func LoginAsTheBotUser() {
	Config.Client.SetToken(Config.Token)

	if user, resp := Config.Client.GetUser(Config.UserId, ""); resp.Error != nil {
		log.Println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		Server.PrintError(resp.Error)
		os.Exit(1)
	} else {
		Config.BotUser = user
	}
}

func UpdateTheBotUserIfNeeded() {
	if Config.BotUser.BotDescription != Config.Description || Config.BotUser.Nickname != Config.DisplayName || Config.BotUser.Username != Config.Username {
		Config.BotUser.BotDescription = Config.Description
		Config.BotUser.Nickname = Config.DisplayName
		Config.BotUser.Username = Config.Username

		if user, resp := Config.Client.UpdateUser(Config.BotUser); resp.Error != nil {
			log.Println("We failed to update the Sample Bot user")
			Server.PrintError(resp.Error)
			os.Exit(1)
		} else {
			Config.BotUser = user
			log.Println("Looks like this might be the first run so we've updated the bots account settings")
		}
	}
}

func FindBotTeam() {
	if team, resp := Config.Client.GetTeamByName(Config.TeamName, ""); resp.Error != nil {
		log.Println("We failed to get the initial load")
		log.Println("or we do not appear to be a member of the team '" + Config.TeamName + "'")
		Server.PrintError(resp.Error)
		os.Exit(1)
	} else {
		Config.BotTeam = team
	}
}

func CreateBotDebuggingChannelIfNeeded() {
	if !Config.Debug {
		return
	}

	if rchannel, resp := Config.Client.GetChannelByName(Config.ChannelLogName, Config.BotTeam.Id, ""); resp.Error != nil {
		log.Println("We failed to get the channels")
		Server.PrintError(resp.Error)
	} else {
		Config.DebuggingChannel = rchannel
		return
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = Config.ChannelLogName
	channel.DisplayName = "Debugging For " + Config.DisplayName
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.CHANNEL_OPEN
	channel.TeamId = Config.BotTeam.Id

	if rchannel, resp := Config.Client.CreateChannel(channel); resp.Error != nil {
		log.Println("We failed to create the channel " + Config.ChannelLogName)
		Server.PrintError(resp.Error)
	} else {
		Config.DebuggingChannel = rchannel
		log.Println("Looks like this might be the first run (with Debug) so we've created the channel " + Config.ChannelLogName)
	}
}
