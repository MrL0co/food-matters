package Config

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/mattermost/mattermost-server/v5/model"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var Client *model.Client4
var WebSocketClient *model.WebSocketClient

var BotUser *model.User
var BotTeam *model.Team
var DebuggingChannel *model.Channel

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var (
	ServerUrl = os.Getenv("MATTERMOST_SERVER")

	Name        = "Food Matters"
	Description = "Food ordering bot"

	DisplayName = getEnv("BOT_DISPLAY_NAME", "Food Matters")

	Username = getEnv("BOT_USERNAME", "food-matters")
	UserId   = os.Getenv("BOT_USER_ID")
	Token    = os.Getenv("BOT_TOKEN")

	CommandToken = os.Getenv("BOT_COMMAND_TOKEN")

	TeamName        = os.Getenv("TEAM")
	FoodChannelName = getEnv("FOOD_CHANNEL", "Order up")

	Debug, _       = strconv.ParseBool(getEnv("DEBUG", "false"))
	ChannelLogName = getEnv("DEBUG_CHANNEL", "debugging-for-food-matters")
)

func GetPort() string {
	return ":8090"
}

func GetServerAddress() string {
	return "http://" + GetOutboundIP().String() + ":8090"
}

func GetWebsocketServer() string {
	return strings.Replace(ServerUrl, "http", "ws", 1)
}

// GetOutboundIP Gets preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
