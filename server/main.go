package main

import (
		"flag"
		"fmt"
		"golang.org/x/net/websocket"
		"log"
		"net/http"
)

type Message struct {
		Name string
		Text string `json:"text"`
}

type hub struct {
		clients 		 map[string] *websocket.Conn
		addClientChan	 chan *websocket.Conn
		removeClientChan chan *websocket.Conn
		broadcastChan	 chan Message
}

var (
		port = flag.String("port", "9000", "port used for ws connection")
)

func server(port string) error {
		hub := newHub()
		mux := http.NewServeMux()
		mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
				handler(ws, hub)
		}))

		s := http.Server{ Addr: ":" + port, Handler: mux}
		return s.ListenAndServe()
}

func handler(ws *websocket.Conn, h *hub)  {
		go	h.run()

		h.addClientChan <- ws

		for {
				var m Message
				err := websocket.JSON.Receive(ws, &m)
				if err != nil {
						h.broadcastChan <- Message{Name: "Server", Text: err.Error()}
						h.removeClient(ws)
						return
				}
				h.broadcastChan <- m
		}
}

func newHub() *hub {
		return &hub{
				clients:		  make(map[string] *websocket.Conn),
				addClientChan:	  make(chan *websocket.Conn),
				removeClientChan: make(chan *websocket.Conn),
				broadcastChan:	  make(chan Message),
		}
}

func (h *hub) run() {
		for {
				select {
				case conn := <- h.addClientChan:
						h.addClient(conn)
				case conn := <- h.removeClientChan:
						h.removeClient(conn)
				case m := <- h.broadcastChan:
						h.broadcastMessage(m)
				}
		}
}

func (h *hub) removeClient(conn *websocket.Conn) {
		delete(h.clients, conn.LocalAddr().String())
}

func (h *hub) addClient(conn *websocket.Conn) {
		h.clients[conn.RemoteAddr().String()] = conn
}

func (h *hub) broadcastMessage(m Message) {
		for _, conn := range h.clients {
				err := websocket.JSON.Send(conn, m)
				if err != nil {
						fmt.Println("Error broadcasting message: ", err)
						return
				}
		}
}

func main()  {
		flag.Parse()
		log.Fatal(server(*port))
}