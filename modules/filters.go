package modules

import "github.com/disgoorg/snowflake/v2"

type MessageFilterFunc func(message string) bool

type MessageFilter struct {
	GuildID snowflake.ID
	Filter  MessageFilterFunc
}
