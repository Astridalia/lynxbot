package commands

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
)

var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Responds with pong!",
}

func HandlePing(e *handler.CommandEvent) error {
	var gatewayPing string
	if e.Client().HasGateway() {
		gatewayPing = e.Client().Gateway().Latency().String()
	}

	eb := discord.NewEmbedBuilder().
		SetTitle("Ping").
		AddField("Rest", "loading...", false).
		AddField("Gateway", gatewayPing, false)

	defer func() {
		var start int64
		_, _ = e.Client().Rest().GetBotApplicationInfo(func(config *rest.RequestConfig) {
			start = time.Now().UnixNano()
		})
		duration := time.Now().UnixNano() - start
		eb.SetField(0, "Rest", time.Duration(duration).String(), false)
		if _, err := e.Client().Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.MessageUpdate{Embeds: &[]discord.Embed{eb.Build()}}); err != nil {
			e.Client().Logger().Error("Failed to update ping embed: ", err)
		}
	}()
	return Respond(e, eb)(nil)
}
