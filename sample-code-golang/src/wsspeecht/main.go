package main

import (
	"flag"
	"log"
	"time"

	pb "pb"
	"wstool"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	host   = flag.String("host", "ws://apigw.open.rokid.com:8444", "Server address")
	lang   = flag.String("lang", "zh", "Language")
	vt     = flag.String("vt", "", "Voice trigger")
	text   = flag.String("text", "", "Text")
	authit = flag.Bool("auth", false, "Need auth?")
	stack  = flag.String("stack", "", "Domain stack")
	device = flag.String("device", "", "Device info")
	count  = flag.Int("count", 1, "Test count")
)

func call_speecht(ws *websocket.Conn, lang, vt, stack, device, text string, authit bool) {
	wspb := wstool.NewWspb(ws)
	id := wstool.RandInt()
	req := &pb.SpeechRequest{
		Id:     proto.Int32(id),
		Type:   pb.ReqType_TEXT.Enum(),
		Lang:   proto.String(lang),
		Vt:     proto.String(vt),
		Stack:  proto.String(stack),
		Device: proto.String(device),
		Asr:    proto.String(text),
	}
	if err := wspb.Write(req); err != nil {
		log.Fatalf("%d write(): %s", id, err)
	}

	start := time.Now()
	for {
		data := &pb.SpeechResponse{}
		if err := wspb.Read(data); err != nil || data.GetResult() != pb.SpeechErrorCode_SUCCESS {
			log.Fatalf("%d read(): %v", id, err)
		}
		log.Printf("Speecht(%s) = asr(%s), nlp(%s), action(%s), cost: %s", text, data.GetAsr(), data.GetNlp(), data.GetAction(), time.Since(start))

		if data.GetFinish() {
			break
		}
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	flag.Parse()

	ws, err := wstool.Dial(*host, wstool.WithAuth(*authit, "speech", "1.0"))
	if err != nil {
		return
	}
	defer ws.Close()

	for i := 0; i < *count; i += 1 {
		call_speecht(ws, *lang, *vt, *stack, *device, *text, *authit)
	}
}
