package main

import "testing"

func TestSpeechv(t *testing.T) {
	conn, err := do_conn("apigw.open.rokid.com:443", true, true)
	if err != nil {
		return
	}
	defer conn.Close()

	call_speechv(conn, "zh", "pcm", "", "", "zhrmghg.pcm", true)
}

func BenchmarkSpeechv(b *testing.B) {
	conn, err := do_conn("apigw.open.rokid.com:443", true, true)
	if err != nil {
		return
	}
	defer conn.Close()

	for i := 0; i < b.N; i++ {
		call_speechv(conn, "zh", "pcm", "", "", "zhrmghg.pcm", true)
	}
}
