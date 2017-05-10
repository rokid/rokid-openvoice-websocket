package wstool

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	pb "pb"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	key            = "rokid_test_key"
	device_type_id = "rokid_test_device_type_id"
	device_id      = "rokid_test_device_id"
	secret         = "rokid_test_secret"
)

type SpeechCredential struct {
	service string
	version string
}

func NewSpeechCredential(service, version string) *SpeechCredential {
	return &SpeechCredential{
		service: service, version: version,
	}
}

func (c SpeechCredential) genAuthMap() map[string]string {
	keys := []string{
		"key",
		"device_type_id",
		"device_id",
		"service",
		"version",
		"time",
	}

	vals := []string{
		key,
		device_type_id,
		device_id,
		c.service,
		c.version,
		strconv.FormatInt(time.Now().Unix(), 10),
	}

	amap := make(map[string]string)
	str := ""
	for n, k := range keys {
		str = str + k + "=" + vals[n] + "&"
		amap[k] = vals[n]
	}
	str = str + "secret=" + secret
	amap["sign"] = fmt.Sprintf("%X", md5.Sum([]byte(str)))

	return amap
}

func (c SpeechCredential) Authws(ws *websocket.Conn) (int32, error) {
	req := &pb.AuthRequest{
		Key:          proto.String(key),
		DeviceTypeId: proto.String(device_type_id),
		DeviceId:     proto.String(device_id),
		Service:      proto.String(c.service),
		Version:      proto.String(c.version),
		Timestamp:    proto.String(strconv.FormatInt(time.Now().Unix(), 10)),
		Sign:         proto.String(c.genAuthMap()["sign"]),
	}

	msg, err := proto.Marshal(req)
	if err != nil {
		log.Printf("marshal(): %s", err)
		return 0, err
	}

	err = ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		log.Printf("write(): %s", err)
		return 0, err
	}

	// read msg
	mt, msg, err := ws.ReadMessage()
	if err != nil || mt != websocket.BinaryMessage {
		log.Printf("read(): %d %s", mt, err)
		return 0, err
	}

	// unmarshal it
	authres := &pb.AuthResponse{}
	err = proto.Unmarshal(msg, authres)
	if err != nil {
		log.Printf("unmarshal(): %s", err)
		return 0, err
	}

	return 0, err
}

func RandInt() int32 {
	seed := time.Now().UTC().UnixNano()
	clientrand := rand.New(rand.NewSource(seed))
	id := clientrand.Int31()

	return id
}

type dialOptions struct {
	service string
	version string
	auth    bool
}

type DialOption func(*dialOptions)

func WithAuth(auth bool, service, version string) DialOption {
	return func(o *dialOptions) {
		o.service = service
		o.version = version
		o.auth = auth
	}
}

func Dial(host string, opts ...DialOption) (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		log.Printf("dial: %s", err)
		return nil, err
	}

	o := &dialOptions{}
	for _, opt := range opts {
		opt(o)
	}

	if o.auth {
		sc := NewSpeechCredential(o.service, o.version)
		if ret, e := sc.Authws(ws); e != nil || ret != 0 {
			log.Printf("client.Auth(): %d %s", ret, e)
			ws.Close()
			return nil, err
		}
	}

	return ws, nil
}
