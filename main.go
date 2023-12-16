package main

import (
	"fmt"
	"os"

	"github.com/astridalia/lynxbot/commands"
	"github.com/astridalia/lynxbot/lynx"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("error while loading .env file: %s", err.Error()))
	}
	token, exists := os.LookupEnv("BOT_TOKEN")
	if !exists {
		panic("BOT_TOKEN environment variable not set")
	}
	bot := lynx.NewBot(token)
	router := handler.New()
	bot.Setup(router)
	ginEngine := gin.New()
	bot.SetupRoutes(ginEngine)

	server_id, exists := os.LookupEnv("SERVER_ID")
	if !exists {
		panic("SERVER_ID environment variable not set")
	}

	bot.SyncCommands(commands.Commands, []snowflake.ID{snowflake.MustParse(server_id)}...)

	RegisterCommandHandlers(router, bot)
	go bot.StartAndBlock()
	err = ginEngine.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("error while running gin engine: %s", err.Error()))
	}
}

func RegisterCommandHandlers(cr handler.Router, b *lynx.Bot) {
	cr.Command("/ping", commands.HandlePing)
}
