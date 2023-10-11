package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/dimeko/wserver/models"

	"golang.org/x/net/websocket"
)

func StartClient(host string, port string, path string, client_name string) {
	var wg sync.WaitGroup

	conn_string := _contruct_full_url(host, port, path)
	log.Printf("Connecting to %s", conn_string)

	ws, err := websocket.Dial(conn_string, "", "http://"+host)
	if err != nil {
		log.Fatal("Connection error, shutting down")
		return
	}
	defer ws.Close()
	wg.Add(1)

	go func() {
		var sending_packet models.ClientToServerMsg
		sending_packet.Broadcast = true
		input := bufio.NewScanner(os.Stdin)
		for {
			input.Scan()

			sending_packet.Message = input.Text()
			sending_packet.Rcv_id = ""
			reqBodyBytes, err := json.Marshal(sending_packet)

			if err != nil {
				log.Println("Could not encode message")
				continue
			}
			_, err = ws.Write(reqBodyBytes)

			if err != nil {
				log.Println("Could not write message")
				continue
			}
		}
	}()

	wg.Add(1)
	go func() {
		var msg = make([]byte, 512)
		var l int
		for {
			l, err = ws.Read(msg)
			if err != nil {
				log.Fatal("Connection with server closed: " + err.Error())
			}
			if l == 0 {
				continue
			}
			fmt.Printf("Incoming message: %s. \n", msg[:l])
		}
	}()

	wg.Wait()
}

func _contruct_full_url(host string, port string, path string) string {
	return fmt.Sprintf("ws://%s:%s/%s", host, port, path)
}
