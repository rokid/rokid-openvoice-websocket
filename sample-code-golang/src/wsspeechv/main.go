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
	vt     = flag.String("vt", "", "Voice trigger")
	fname  = flag.String("file", "", "Audio file")
	authit = flag.Bool("auth", false, "Need auth?")
	stack  = flag.String("stack", "", "Current domain")
	device = flag.String("device", "", "Device info")
	count  = flag.Int("count", 1, "Test count")
)

func call_speechv(ws *websocket.Conn, lang, codec, stack, device, vt, fname string, authit bool) {
	var file *os.File
	if f, err := os.Open(fname); err != nil {
		log.Fatalf("could not open file %v: %v", fname, err)
	} else {
		file = f
	}
	defer file.Close()

	wspb := wstool.NewWspb(ws)
	id := wstool.RandInt()

	start := time.Now()
	waitc := make(chan struct{})
	go func() {
		for {
			data := &pb.SpeechResponse{}
			if err := wspb.Read(data); err != nil || data.GetResult() != pb.SpeechErrorCode_SUCCESS {
				log.Fatalf("%d read(): %v", id, err)
			}
			log.Printf("Speechv(%s) = asr(%s), nlp(%s), action(%s), cost: %s", fname, data.GetAsr(), data.GetNlp(), data.GetAction(), time.Since(start))

			if data.GetFinish() {
				break
			}
		}
		close(waitc)
	}()

	req := &pb.SpeechRequest{
		Id:     proto.Int32(id),
		Type:   pb.ReqType_START.Enum(),
		Lang:   proto.String(lang),
		Vt:     proto.String(vt),
		Codec:  proto.String(codec),
		Stack:  proto.String(stack),
		Device: proto.String(device),
	}
	if err := wspb.Write(req); err != nil {
		log.Fatalf("%d write(): %s", id, err)
	}

	voice := make([]byte, 320*30)
	for {
		time.Sleep(300 * time.Millisecond)

		if n, err := file.Read(voice[:]); err == nil {
			log.Printf("Read file(%d)", n)
			req := &pb.SpeechRequest{
				Id:    proto.Int32(id),
				Type:  pb.ReqType_VOICE.Enum(),
				Voice: voice[:n],
			}
			if err := wspb.Write(req); err != nil {
				log.Fatalf("%d write(): %s", id, err)
			}
		} else {
			log.Printf("Read file: %v", err)
			break
		}
	}
	req = &pb.SpeechRequest{
		Id:   proto.Int32(id),
		Type: pb.ReqType_END.Enum(),
	}
	if err := wspb.Write(req); err != nil {
		log.Fatalf("%d write(): %s", id, err)
	}

	<-waitc
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
		call_speechv(ws, *lang, *codec, *stack, *device, *vt, *fname, *authit)
	}
}
