target := wstts wsasr wsspeechv wsspeecht

all: pre $(target)

pre: src/pb/speech.pb.go

src/pb/speech.pb.go: src/pb/speech.proto
	cd src/pb; protoc --go_out=plugins=grpc:. *.proto
	go get $(target)

wstts: pre
	go build $@

wsasr: pre
	go build $@

wsspeechv: pre
	go build $@

wsspeecht: pre
	go build $@

clean:
	rm -vf $(target) src/pb/speech.pb.go
