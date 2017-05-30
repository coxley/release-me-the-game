package main

import (
	"flag"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	// "github.com/coxley/release-me-the-game/state"
)

var addr = flag.String("addr", "localhost:8080", "http server address")

func main() {
	flag.Parse()
	glog.Info("Attempting to connect to server")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}

	glog.Infof("Connecting to %s", u)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Fatalf("Failed to connect to server: %s", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				glog.Fatalf("Read: %s", err)
				return
			}
			glog.Infof("Recv: %s", msg)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			char := state.Character{}
			data, _ := proto.Marshal(char)
			err := c.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				glog.Warningf("Write: %s", err)
				return
			}
		case <-interrupt:
			glog.Info("Interrupt")
			// To cleanly close, client should send a close frame and wait for
			// server
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				glog.Warningf("Write close: %s", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}
