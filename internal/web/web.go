package web

import (
	"fmt"
	"net/http"

	"github.com/Oxel40/hermes/internal/configuration"
	"github.com/Oxel40/hermes/internal/logging"
	"github.com/Oxel40/hermes/internal/token"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type toSend struct {
	From string   `json:"from"`
	Msg  string   `json:"msg"`
	IDs  []string `json:"ids"`
}

type Web struct {
	communicatorTokenMap *token.TokenMap
	serviceTokenMap      *token.TokenMap
	log                  *logging.Logger
	config               *configuration.Config
	nameToConn           map[string]*websocket.Conn
}

// AttatchServiceTokenMap ...
func (web *Web) AttatchServiceTokenMap(tokenMap *token.TokenMap) {
	web.serviceTokenMap = tokenMap
}

// AttatchCommunicatorTokenMap ...
func (web *Web) AttatchCommunicatorTokenMap(tokenMap *token.TokenMap) {
	web.communicatorTokenMap = tokenMap
}

// AttatchLogger ...
func (web *Web) AttatchLogger(log *logging.Logger) {
	web.log = log
}

// AttatchConfig ...
func (web *Web) AttatchConfig(config *configuration.Config) {
	web.config = config
}

func (web *Web) Subroutine(port int) {
	web.nameToConn = make(map[string]*websocket.Conn)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	router := mux.NewRouter()

	router.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	router.HandleFunc("/communicator", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusUpgradeRequired)
			w.Header().Set("Upgrade", "WS")
			fmt.Fprint(w, "Upgrade required\n")
			web.log.Warning.Println(r.RemoteAddr, "couldn't upgrade connection:", err)
			return
		}

		// Read authentication token from client
		_, token, err := conn.ReadMessage()
		if err != nil {
			return
		}
		name := web.communicatorTokenMap.TokenToName[string(token)]
		if name == "" {
			// Invalid token
			// Write message back to client
			conn.WriteMessage(websocket.CloseMessage, []byte("Invalid token"))
			conn.Close()
			return
		}

		if web.nameToConn[name] != nil {
			web.nameToConn[name].WriteMessage(websocket.CloseMessage, []byte(r.RemoteAddr+" has taken over this connection"))
			web.nameToConn[name].Close()
			web.log.Warning.Println(r.RemoteAddr, "took over the", name, "connection from", err)
		}
		web.nameToConn[name] = conn

		conn.SetCloseHandler(func(code int, text string) error {
			web.log.Trace.Println(name, "disconnected:", code, text)
			web.nameToConn[name] = nil
			return nil
		})

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	router.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		msg := r.FormValue("msg")
		name := web.serviceTokenMap.TokenToName[token]

		web.log.Trace.Printf("%s tries to authenticate with token: %s\n", r.RemoteAddr, token)
		if name == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid token\n")
			web.log.Trace.Printf("%s failed, invalid token\n", r.RemoteAddr)
			return
		}
		web.log.Trace.Printf("%s is %s, sending message\n", r.RemoteAddr, name)

		go brodcastMessage(web, name, msg)

		fmt.Fprintf(w, "Message \"%s\" by %s has been passed on.\n", msg, name)
	}).Methods("POST")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome the hermes REST API! Please consult https://github.com/Oxel40/hermes for API specific information")
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func brodcastMessage(web *Web, from, msg string) {
	subs := getSubscribers(web.config, from)

	for _, dis := range web.config.Communicators {
		go func(idIndex int, distribName string) {
			conn := web.nameToConn[distribName]
			if conn != nil {
				var out toSend
				out.From = from
				out.Msg = msg
				out.IDs = make([]string, 0)
				for _, sub := range subs {
					out.IDs = append(out.IDs, sub.IDs[idIndex])
				}
				err := conn.WriteJSON(out)
				if err != nil {
					web.log.Error.Println("Couldn't send to:", distribName, ":", err)
				}
			}
		}(dis.IDIndex, dis.Name)
	}
}

func getSubscribers(config *configuration.Config, serviceName string) []configuration.Recipient {
	out := make([]configuration.Recipient, 0)
	for _, rep := range config.Recipiens {
		for _, serv := range rep.Subscriptions {
			if serviceName == serv {
				out = append(out, rep)
				break
			}
		}
	}
	return out
}
