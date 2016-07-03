package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"net/http"
	"net/url"
)

var upgrader = websocket.Upgrader{}

// websocket handler
func startTunnel(rw http.ResponseWriter, req *http.Request) {

	wsoconn, err := upgrader.Upgrade(rw, req, nil)

	if err != nil {
		fmt.Println("Failed to upgrade http request. Developer details: ", err.Error(), "\nExiting")
		return
	}

	// parse path from request
	query := req.URL.Query()

	path := query["path"][0]

	// set up tail
	tail, err := tail.TailFile(path, tail.Config{
		Follow: true,
		ReOpen: true,
		Poll:   true,
	})

	if err != nil {
		fmt.Println("Failed to track file ",
			path,
			"\nPath is either invalid, or has restricted permissions. Developer details: ",
			err.Error())
		return
	}

	// tail file and send new lines over websocket
	for line := range tail.Lines {

		err = wsoconn.WriteMessage(8398, []byte(line.Text))

		if err != nil {
			fmt.Println("Failed to send message through tunnel. Tunnel must have been broken. Developer details: ",
				err.Error())
			return
		}
	}
}
