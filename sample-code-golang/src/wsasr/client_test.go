package main

import "testing"

func TestAsr(t *testing.T) {
	conn, err := do_conn("apigw.open.rokid.com:443", true, true)
	if err != nil {
		return
	}
	defer conn.Close()

	call_asr(conn, "zh", "pcm", "zhrmghg.pcm", true)
}

func BenchmarkAsr(b *testing.B) {
	conn, err := do_conn("apigw.open.rokid.com:443", true, true)
	if err != nil {
		return
	}
	defer conn.Close()

	for i := 0; i < b.N; i++ {
		call_asr(conn, "zh", "pcm", "zhrmghg.pcm", true)
	}
}
