package lynx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/gin-gonic/gin"
)

type Bot struct {
	Token  string
	Client bot.Client
	Mux    *http.ServeMux
}

func NewBot(token string) *Bot {
	return &Bot{
		Token: token,
	}
}

func (b *Bot) SetupRoutes(engine *gin.Engine) {
	engine.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, Gin!")
	})
}

func (b *Bot) clientConfigurator(r handler.Router) []bot.ConfigOpt {
	return []bot.ConfigOpt{
		bot.WithEventListeners(r),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildMessages|
				gateway.IntentDirectMessages|
				gateway.IntentGuildMessageTyping|
				gateway.IntentDirectMessageTyping,
			),
			gateway.WithCompress(true),
		),

		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagGuilds)),
	}
}

func (b *Bot) Setup(r handler.Router) {
	var err error
	b.Client, err = disgo.New(b.Token, b.clientConfigurator(r)...)
	if err != nil {
		log.Fatalf("error while building disgo client: %s", err.Error())
	}

	if err = b.Client.OpenGateway(context.Background()); err != nil {
		log.Fatalf("error while opening gateway connection: %s", err.Error())
	}
}

func (b *Bot) StartAndBlock() {
	log.Println("Bot is running... Press Ctrl-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s

	defer b.Shutdown()
}

func (b *Bot) Shutdown() {
	log.Println("Shutting down bot...")
	b.Client.Close(context.Background())
	os.Exit(0)
}

func (b *Bot) SyncCommands(commands []discord.ApplicationCommandCreate, guildIDs ...snowflake.ID) {
	restClient := b.Client.Rest()
	appID := b.Client.ApplicationID()

	// Common function to set commands
	setCommands := func(id snowflake.ID, global bool) {
		var err error
		if global {
			_, err = restClient.SetGlobalCommands(appID, commands)
		} else {
			_, err = restClient.SetGuildCommands(appID, id, commands)
		}
		if err != nil {
			log.Fatalf("failed to set commands%s: %s", func() string {
				if global {
					return " globally"
				}
				return fmt.Sprintf(" for guild %s", id)
			}(), err)
		}
	}

	// Determine whether to set global or guild-specific commands
	if len(guildIDs) == 0 {
		setCommands(0, true) // Set global commands
	} else {
		for _, guildID := range guildIDs {
			setCommands(guildID, false) // Set guild-specific commands
		}
	}
}
