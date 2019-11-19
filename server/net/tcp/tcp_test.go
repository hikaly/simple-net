package tcp

import (
	"encoding/binary"
	"testing"
)

func Test_DataFilter(t *testing.T) {
	msg_src := []byte("this is a test message.")
	msg_h := make([]byte, 2)
	binary.BigEndian.PutUint16(msg_h, uint16(len(msg_src)))

	// 完整的
	msg_b := append(msg_h, msg_src...)
	DataFilter("1", msg_b, len(msg_b))

	// 超出的
	msg_b2 := append(msg_h, msg_src...)
	msg_b2 = append(msg_b2, msg_h...)
	DataFilter("1", msg_b2, len(msg_b2))

	// 剩余的
	msg_b3 := []byte{}
	msg_b3 = append(msg_b3, msg_src...)
	DataFilter("1", msg_b3, len(msg_b3))

	msg_b4 := []byte{}
	msg_b4 = append(msg_b4, msg_h...)
	msg_b4 = append(msg_b4, msg_src...)
	msg_b4 = append(msg_b4, msg_h...)
	DataFilter("1", msg_b4, len(msg_b4))

	msg_b5 := []byte{}
	msg_b5 = append(msg_b5, msg_src...)
	DataFilter("1", msg_b5, len(msg_b5))
}

func Benchmark_DataFilter(b *testing.B) {
	msg_src := []byte{}
	for {
		msg_src = append(msg_src, []byte("this is a test message.")...)

		if len(msg_src) > 32*1024 {
			break
		}
	}

	msg_h := make([]byte, 2)
	binary.BigEndian.PutUint16(msg_h, uint16(len(msg_src)))

	msg_b := append(msg_h, msg_src...)

	for i := 0; i < b.N; i++ {
		DataFilter("1", msg_b, len(msg_b))
	}
}
