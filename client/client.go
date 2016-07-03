package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var signals = make(chan os.Signal, 1)
var wsConn *websocket.Conn

func main() {

	// verify command line argument (usage ./client http://url filepath)
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./tunnel-client destination filepath")
		os.Exit(1)
	}

	// set up signal handling for graceful exit
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTSTP)

	go func() {
		_ = <-signals
		fmt.Println("Terminating...")
		wsConn.Close()
		os.Exit(0)

	}()

	// get destination and filepath
	dest, path := os.Args[1], os.Args[2]

	// create url
	url := url.URL{
		Scheme: "ws",
		Host:   dest,
		Path:   "/tunnel",
	}

	// set query values
	query := url.Query()
	query.Set("path", path)
	url.RawQuery = query.Encode()

	// open websocket
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	// allow signal handler to access websocket conn for graceful exit
	wsConn = conn

	if err != nil {
		fmt.Println("Failed to establish tunnel\nDeveloper Details: ", err.Error())
		os.Exit(2)
	}

	fmt.Println("Tunnel Established. Now streaming ", path, " from: ", url.String(), "\n\n")

	for {

		// start streaming to command line
		_, log, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Tunnel closed\nDeveloper Details: ", err.Error(), "\nExiting...")
			os.Exit(4)
		}

		// print message to screen
		fmt.Println(string(log))
	}

}
