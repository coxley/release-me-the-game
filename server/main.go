package main

import (
	"flag"
	"net/http"
	// "fmt"

	// "net/http"
	"github.com/coxley/release-me-the-game/types"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http server address")
var upgrader = websocket.Upgrader{}

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	http.HandleFunc("/echo", echo)

	glog.Infof("Listening on %s ...", *addr)
	glog.Fatal(http.ListenAndServe(*addr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Errorf("Failed to upgrade client connection to WS: %s", err)
	}
	defer c.Close()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			glog.V(3).Infof("Read: %s", err)
			break
		}
		char := &types.Character{}
		_ = proto.Unmarshal(msg, char)
		glog.V(3).Infof("Recv: %#v", char)
		err = c.WriteMessage(mt, msg)
		if err != nil {
			glog.V(3).Infof("Write: %s", err)
			break
		}
	}
}
