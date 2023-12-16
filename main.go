package main

import (
	"fmt"
	"os"

	"github.com/astridalia/lynxbot/lynx"
	"github.com/disgoorg/disgo/handler"
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
	go bot.StartAndBlock()
	err = ginEngine.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("error while running gin engine: %s", err.Error()))
	}

}
