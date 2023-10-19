package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverURL := "ws://192.168.10.197:8080/ws"
	u, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal(err)
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ticker.C:
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			message := []byte(input)
			err := c.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Send error:", err)
				return
			}
			fmt.Printf("Sent message: %s\n", message)

			_, response, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Receive error:", err)
				return
			}
			fmt.Printf("Received message: %s\n", response)

		case <-interrupt:
			fmt.Println("Received interrupt signal, closing connection...")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("Close connection error:", err)
				return
			}
			select {
			case <-time.After(1 * time.Second):
			}
			return
		}
	}
}
