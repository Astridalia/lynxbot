package commands

import (
	"encoding/json"

	"github.com/astridalia/lynxbot/mediawiki"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type WikiContext struct {
	Name        string `json:"name"`
	School      string `json:"school"`
	Description string `json:"descrip1"`
	EggType     string `json:"egg"`
	Accuracy    string `json:"accuracy"`
	PvPlevel    string `json:"PvPlevel"`
	Image       string `json:"image"`
}

var wikiCommand = discord.SlashCommandCreate{
	Name:        "wiki",
	Description: "Searches the wiki for a given term",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "term",
			Description: "The term to search for",
			Required:    true,
		},
	},
}

func HandleWiki(e *handler.CommandEvent) error {
	term := e.SlashCommandInteractionData().String("term")
	eb := discord.NewEmbedBuilder()
	wiki := mediawiki.NewWikiService()
	text, err := wiki.Json("TreasureCard:" + term)
	if err != nil {
		eb.SetDescription("Error: " + err.Error())
		return Respond(e, eb)
	}
	var ctx WikiContext
	err = json.Unmarshal(text, &ctx)
	if err != nil {
		HandleError(eb, err)
		return Respond(e, eb)
	}

	eb.SetTitle(ctx.Name)
	eb.SetImage(ctx.Image)
	eb.AddFields(
		discord.EmbedField{
			Name:  "School",
			Value: ctx.School,
		},
		discord.EmbedField{
			Name:  "Description",
			Value: ctx.Description,
		},
		discord.EmbedField{
			Name:  "Accuracy",
			Value: ctx.Accuracy,
		},
		discord.EmbedField{
			Name:  "PvP Level",
			Value: ctx.PvPlevel,
		},
	)

	return Respond(e, eb)
}

func Respond(e *handler.CommandEvent, eb *discord.EmbedBuilder) error {
	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().SetEmbeds(eb.Build()).SetEphemeral(true).Build(),
	)
}

func HandleError(eb *discord.EmbedBuilder, err error) *discord.EmbedBuilder {
	eb.SetDescription("Error: " + err.Error())
	eb.SetColor(0xFF0000)
	return eb
}
