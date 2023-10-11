package server

import (
	"encoding/json"
	"log"

	"github.com/dimeko/wserver/models"
)

func HubOrchestrator(server *Server) {
	for {
		select {
		case incmng_msg := <-server.message_orch:

			outgoing_message, err := json.Marshal(&models.ServerToClientMsg{
				Sndr_id:   incmng_msg.Sndr_id.String(),
				Message:   incmng_msg.Message,
				Broadcast: incmng_msg.Broadcast,
			})
			if err != nil {
				log.Printf("Error encoding outgoing message: %s", err.Error())
				continue
			}
			if client, ok := server.conns[incmng_msg.Rcv_id]; ok {
				log.Printf("Sending message to client: %s", incmng_msg.Rcv_id)
				_, err := client.conn.Write([]byte(outgoing_message))
				if err != nil {
					log.Printf("Error sending message: %s", err.Error())
				}
			} else {
				if incmng_msg.Broadcast {
					for _, client := range server.conns {
						_, err := client.conn.Write([]byte(outgoing_message))
						if err != nil {
							log.Printf("Error sending message: %s", err.Error())
						}
					}
				} else {
					log.Println("Wrong user UUID sent.")
				}
			}
		}
	}
}
