package main

import "testing"

func TestTts(t *testing.T) {
	conn, err := do_conn("apigw.open.rokid.com:443", true, true)
	if err != nil {
		return
	}
	defer conn.Close()

	call_tts(conn, "zh", "opu", "map 是 Golang 中的一种 Associative data type。提供类似于其他语言中 hash 或者 dictionary 的功能 。", "", true)
}

func BenchmarkTts(b *testing.B) {
	conn, err := do_conn("apigw.open.rokid.com:443", true, true)
	if err != nil {
		return
	}
	defer conn.Close()

	for i := 0; i < b.N; i++ {
		call_tts(conn, "zh", "opu", "map 是 Golang 中的一种 Associative data type。提供类似于其他语言中 hash 或者 dictionary 的功能 。", "", true)
	}
}
