package modules

import "github.com/disgoorg/snowflake/v2"

type MessageFilter struct {
	GuildID snowflake.ID
	Filter  func(message string) bool
}
