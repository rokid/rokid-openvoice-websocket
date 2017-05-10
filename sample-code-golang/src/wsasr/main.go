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
	codec  = flag.String("codec", "pcm", "Codec")
	fname  = flag.String("file", "", "Audio file")
	authit = flag.Bool("auth", false, "Need auth?")
	tls    = flag.Bool("tls", false, "Need tls?")
	count  = flag.Int("count", 1, "Test count")
)

func call_asr(ws *websocket.Conn, lang, codec, fname string, authit bool) {
	var file *os.File

	if f, err := os.Open(fname); err != nil {
		log.Fatalf("could not open file %v: %v", fname, err)
	} else {
		file = f
	}
	defer file.Close()

	id := wstool.RandInt()
	req := &pb.AsrRequest{
		Id:    proto.Int32(id),
		Type:  pb.ReqType_START.Enum(),
		Lang:  proto.String(lang),
		Codec: proto.String(codec),
	}
	wspb := wstool.NewWspb(ws)
	if err := wspb.Write(req); err != nil {
		log.Fatalf("%d write(): %s", id, err)
	}

	voice := make([]byte, 320*30)
	for {
		var n int
		var err error

		time.Sleep(300 * time.Millisecond)
		isend := false
		if n, err = file.Read(voice[:]); err == nil {
			log.Printf("%d Read file(%d)", id, n)
		} else {
			log.Printf("%d Read file: %v", id, err)
			isend = true
		}
		req := &pb.AsrRequest{
			Id:    proto.Int32(id),
			Type:  pb.ReqType_VOICE.Enum(),
			Lang:  proto.String(lang),
			Codec: proto.String(codec),
			Voice: voice[:n],
		}
		if isend {
			req.Type = pb.ReqType_END.Enum()
		}

		if err = wspb.Write(req); err != nil {
			log.Fatalf("%d write(): %s", id, err)
		}

		if isend {
			break
		}
	}

	start := time.Now()
	waitc := make(chan struct{})
	go func() {
		for {
			data := &pb.AsrResponse{}
			if err := wspb.Read(data); err != nil || data.GetResult() != pb.SpeechErrorCode_SUCCESS {
				log.Fatalf("%d read(): %v", id, err)
			}
			log.Printf("%d Got asr: Finish(%t), Asr('%s'), cost: %s", id, data.GetFinish(), data.GetAsr(), time.Since(start))

			if data.GetFinish() {
				close(waitc)
				break
			}
		}
	}()

	<-waitc
}

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	flag.Parse()

	ws, err := wstool.Dial(*host, wstool.WithAuth(*authit, "asr", "1.0"))
	if err != nil {
		return
	}
	defer ws.Close()

	for i := 0; i < *count; i += 1 {
		call_asr(ws, *lang, *codec, *fname, *authit)
	}
}
