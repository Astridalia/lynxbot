package modules

import "github.com/disgoorg/snowflake/v2"

type EmojiReaction struct {
	EmojiId snowflake.ID `json:"emoji_id"`
	RoleId  snowflake.ID `json:"role_id"`
}

type MessageReaction struct {
	MessageId snowflake.ID    `json:"message_id"`
	Reactions []EmojiReaction `json:"reactions"`
}
