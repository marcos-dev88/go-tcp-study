package tcp_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func Handle(conn net.Conn) {

	for {
		decoder := json.NewDecoder(conn)

		msg := struct {
			Message string `json:"message"`
		}{}

		err := decoder.Decode(&msg)

		if err != nil {
			log.Printf("error: %v", err)
		}

		fmt.Printf("%s \n", msg.Message)
	}
}
