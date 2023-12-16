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

func WikiText(e *handler.CommandEvent, term string) []byte {
	wiki := mediawiki.NewWikiService()
	text, err := wiki.Json("TreasureCard:" + term)
	if err != nil {
		return nil
	}
	return text
}

func WikiContent(e *handler.CommandEvent, term string) (WikiContext, error) {
	var ctx WikiContext
	err := json.Unmarshal(WikiText(e, term), &ctx)
	if err != nil {
		return WikiContext{}, err
	}
	return ctx, nil
}

func BuildWikiEmbed(e *handler.CommandEvent, ctx WikiContext) *discord.EmbedBuilder {
	eb := discord.NewEmbedBuilder()
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
	return eb
}

func HandleWiki(e *handler.CommandEvent) error {
	term := e.SlashCommandInteractionData().String("term")
	context, err := WikiContent(e, term)
	if err != nil {
		return Respond(e, HandleError(err))
	}
	return Respond(e, BuildWikiEmbed(e, context))
}

func Respond(e *handler.CommandEvent, eb *discord.EmbedBuilder) error {
	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().SetEmbeds(eb.Build()).SetEphemeral(true).Build(),
	)
}

func HandleError(err error) *discord.EmbedBuilder {
	eb := discord.NewEmbedBuilder()
	eb.SetDescription("Error: " + err.Error())
	eb.SetColor(0xFF0000)
	return eb
}
