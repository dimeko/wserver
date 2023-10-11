package models

type ClientToServerMsg struct {
	Rcv_id  string `json:"rcv_id"`
	Message string `json:"message"`
	// Room      string `json:"room"`
	Broadcast bool `json:"broadcast"`
}

type ServerToClientMsg struct {
	Sndr_id string `json:"sndr_id"`
	Message string `json:"message"`
	// Room      string `json:"room"`
	Broadcast bool `json:"broadcast"`
}
