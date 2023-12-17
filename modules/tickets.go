package modules

type Ticket struct {
	TicketId   int64  `json:"ticket_id"`
	ChannelId  int64  `json:"channel_id"`
	UserId     int64  `json:"user_id"`
	AssignedBy int64  `json:"assigned_by"`
	Expiration int64  `json:"expiration"`
	Open       bool   `json:"open"`
	Reason     string `json:"reason"`
}

type TicketList struct {
	Tickets []Ticket `json:"tickets"`
}
