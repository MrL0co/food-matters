package Server

import (
	"food-matters/Config"
	"github.com/mattermost/mattermost-server/v5/model"
	"log"
	"os"
)

var (
	Channels = make(map[string]*model.Channel)
)

func RemoveAllMessagesWithText(msg string) {
	if list, resp := Config.Client.SearchPosts(Config.BotTeam.Id, msg, false); resp.Error != nil {
		log.Println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		for _, post := range list.Posts {
			if post.UserId == Config.BotUser.Id {
				if success, resp := Config.Client.DeletePost(post.Id); resp.Error != nil || !success {
					log.Printf("failed to remove post %s, %s", post.Message, post.Id)
				} else {
					log.Printf("removed old post: %s, %s", post.Message, post.Id)
				}
			}
		}
	}
}

func SendMsgToDebuggingChannel(msg string, replyToId string) {
	if !Config.Debug {
		log.Printf("%s - '%s'", replyToId, msg)
		return
	}

	SendMsgToChannel(Config.DebuggingChannel.Id, msg, replyToId)
}

func SendMsgToChannel(channelId string, msg string, replyToId string) {
	channel := GetChannel(channelId)

	post := &model.Post{}
	post.ChannelId = channelId
	post.Message = msg

	post.RootId = replyToId

	if _, resp := Config.Client.CreatePost(post); resp.Error != nil {
		log.Printf("We failed to send a message to %s", channel.DisplayName)
		PrintError(resp.Error)
	}

	log.Printf("Send message to %s", channel.DisplayName)
}

func GetChannel(channelId string) *model.Channel {
	if val, ok := Channels[channelId]; ok {
		return val
	}

	channel, resp := Config.Client.GetChannel(channelId, "")
	if resp.Error != nil {
		log.Fatalln(resp.Error.Error())
	}

	Channels[channelId] = channel
	return channel
}
