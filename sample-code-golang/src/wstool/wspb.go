package wstool

import (
	"fmt"
	"log"
	"syscall"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Wspb struct {
	ws *websocket.Conn
}

func NewWspb(ws *websocket.Conn) *Wspb {
	return &Wspb{
		ws: ws,
	}
}

func (w Wspb) Read(pb proto.Message) error {
	// read msg
	mt, msg, err := w.ws.ReadMessage()
	if err != nil || mt != websocket.BinaryMessage {
		log.Printf("unknown message type %d, %s", mt, err)
		return fmt.Errorf("unknown message type %d, %s", mt, err)
	}

	// unmarshal it
	err = proto.Unmarshal(msg, pb)
	if err != nil {
		log.Printf("unmarshal() %s", err)
		return fmt.Errorf("unmarshal(): %s", err)
	}

	return nil
}

func (w Wspb) Write(pb proto.Message) error {
	msg, err := proto.Marshal(pb)
	if err != nil {
		log.Printf("marshal(): %v", err)
		return err
	}

	if err = w.ws.WriteMessage(websocket.BinaryMessage, msg); err != nil {
		log.Printf("write(): %v", err)
		if err == syscall.EPIPE {
			log.Printf("write(): %v, skip it", err)
		} else {
			return err
		}
	}

	return nil
}
