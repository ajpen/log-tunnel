package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"os"
)

func main() {

	// verify command line argument (usage ./client http://url filepath)
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./tunnel-client destination filepath")
		os.Exit(1)
	}

	dest, path := os.Args[1], os.Args[2]

	// build URL
	url := fmt.Sprintf("%s?path=%s", dest, path)

	// open websocket
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		fmt.Println("Failed to establish tunnel\nDeveloper Details: ", err.Error())
		os.Exit(2)
	}

	fmt.Println("Tunnel Established. Now streaming ", path, " from ", dest, "\n\n")

	for {

		// start streaming to command line
		_, log, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Tunnel closed\nDeveloper Details: ", err.Error(), "\nExiting...")
			os.Exit(4)
		}

		// print message to screen
		fmt.Printf("%s", string(log))
	}

}
