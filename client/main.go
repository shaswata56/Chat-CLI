package main

import (
		"bufio"
		"flag"
		"fmt"
		"golang.org/x/net/websocket"
		"log"
		"math/rand"
		"os"
		"time"
)

type Message struct {
		Name string
		Text string `json:"text"`
}
var (
		port = flag.String("port", "9000", "port used for ws connection")
)

func connect() (*websocket.Conn, error) {
		return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", mockedIP())
}

func mockedIP() string {
		var arr [4]int
		for i := 0; i < 4; i++ {
				rand.Seed(time.Now().UnixNano())
				arr[i] = rand.Intn(256)
		}
		return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}

func main()  {
		flag.Parse()

		name := ""
		fmt.Print("Username: ")
		fmt.Scan(&name)

		ws, err := connect()
		if err != nil {
				log.Fatal(err)
		}
		defer ws.Close()

		var m Message
		go func() {
				for {
						err := websocket.JSON.Receive(ws, &m)
						if err != nil {
								fmt.Println("Error receiving message: ", err.Error())
								break
						}
						if m.Name != name {
								fmt.Println(m.Name ," : ", m.Text)
						}
				}
		}()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
				text := scanner.Text()
				if text == "" {
						continue
				}
				m := Message{
						Name: name,
						Text: text,
				}
				err = websocket.JSON.Send(ws, m)
				if err != nil {
						fmt.Println("Error sending message: ", err.Error())
						break
				}
		}
}