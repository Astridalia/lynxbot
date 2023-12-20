package commands

import (
	"fmt"
	"math/rand"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var randomizerCommand = discord.SlashCommandCreate{
	Name:        "random",
	Description: "responds with a random user",
}

func HandleRandomizer(e *handler.CommandEvent) error {
	if e == nil || e.GuildID() == nil || e.Client() == nil || e.Client().Rest() == nil {
		return fmt.Errorf("nil value detected")
	}

	members, err := e.Client().Rest().GetMembers(*e.GuildID(), 325, discord.AllGuildChannels(*e.GuildID()))
	if err != nil {
		return err
	}

	eb := discord.NewEmbedBuilder()

	if len(members) == 0 {
		return Respond(e, eb.SetDescription("No members found"))(nil)
	}

	// pick a random member
	member := members[rand.Intn(len(members))]

	eb.SetTitle("Winner").
		SetDescription(member.User.Username)

	return Respond(e, eb)(nil)
}
