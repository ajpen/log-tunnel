package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{}

// websocket handler
func startTunnel(rw http.ResponseWriter, req *http.Request) {

	wsoconn, err := upgrader.Upgrade(rw, req, nil)

	if err != nil {
		fmt.Println("Failed to upgrade http request. Developer details: ", err.Error(), "\n")
		return
	}

	defer wsoconn.Close()

	// parse path from request
	query := req.URL.Query()

	var path string

	if p, ok := query["path"]; ok {
		if len(p) > 0 {
			path = p[0]

		} else {
			fmt.Println("Failed to specify path in request.")
			return
		}

	} else {
		fmt.Println("Failed to specify path in request.")
		return
	}

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

		err = wsoconn.WriteMessage(websocket.TextMessage, []byte(line.Text))

		if err != nil {
			fmt.Println("Failed to send message through tunnel. Tunnel must have been broken. Developer details: ",
				err.Error())
			return
		}
	}
}

func main() {

	if len(os.Args) != 3 {
		panic("Usage: ./tunnel-server address port")
	}

	address, port := os.Args[1], os.Args[2]

	listenOn := fmt.Sprintf("%s:%s", address, port)

	http.HandleFunc("/tunnel", startTunnel)
	if err := http.ListenAndServe(listenOn, nil); err != nil {
		panic(err)
	}
}
