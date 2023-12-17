package modules

import "github.com/disgoorg/snowflake/v2"

type SmartRole struct {
	RoleId      snowflake.ID `json:"role_id"`
	Name        string       `json:"name"`
	Permissions int64        `json:"permissions"`
}

type SmartRoleList struct {
	Roles []SmartRole `json:"roles"`
}

type AssignUserRole struct {
	UserId     snowflake.ID `json:"user_id"`
	RoleId     snowflake.ID `json:"role_id"`
	AssignedBy snowflake.ID `json:"assigned_by"`
	Expiration int64        `json:"expiration"`
}
