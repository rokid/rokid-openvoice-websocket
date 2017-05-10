package main

import (
	"flag"
	"log"
	"os"
	"time"

	pb "pb"
	"wstool"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	host   = flag.String("host", "ws://apigw.open.rokid.com:8444", "Server address")
	lang   = flag.String("lang", "zh", "Language")
	codec  = flag.String("codec", "opu2", "Codec")
	text   = flag.String("text", "今天天气怎么样?", "Tts Text")
	fname  = flag.String("file", "", "Out file")
	authit = flag.Bool("auth", false, "Need auth?")
	count  = flag.Int("count", 1, "Test count")
)

func call_tts(ws *websocket.Conn, lang, codec, text, fname string, authit bool) {
	var file *os.File

	if 0 != len(fname) {
		if f, err := os.Create(fname); err != nil {
			log.Fatalf("could not create file %v: %v", fname, err)
		} else {
			file = f
		}
		defer file.Close()
	}

	id := wstool.RandInt()
	req := &pb.TtsRequest{
		Id:        proto.Int32(id),
		Declaimer: proto.String(lang),
		Codec:     proto.String(codec),
		Text:      proto.String(text),
	}
	wspb := wstool.NewWspb(ws)
	if err := wspb.Write(req); err != nil {
		log.Fatalf("%d write(): %s", id, err)
	}

	start := time.Now()
	for {
		data := &pb.TtsResponse{}
		if err := wspb.Read(data); err != nil || data.GetResult() != pb.SpeechErrorCode_SUCCESS {
			log.Fatalf("%d read(): %v", id, err)
		}
		log.Printf("%d Got tts: Finish(%t), Text('%s'), Voice(len=%d), cost: %s", id, data.GetFinish(), data.GetText(), len(data.Voice), time.Since(start))

		if file != nil {
			file.Write(data.Voice)
		}

		if data.GetFinish() {
			break
		}
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	flag.Parse()

	ws, err := wstool.Dial(*host, wstool.WithAuth(*authit, "tts", "1.0"))
	if err != nil {
		return
	}
	defer ws.Close()

	for i := 0; i < *count; i += 1 {
		call_tts(ws, *lang, *codec, *text, *fname, *authit)
	}
}
