package tcp

import (
	"encoding/binary"
)

type MessagePack struct {
	// TCP链接id
	Id string

	// 消息长度
	Size int

	// 消息数据
	Data []byte
}

type Uncomplete struct {
	// 消息包
	Message *MessagePack

	// 已读长度
	Read int

	// 包头
	Header []byte
}

//
var uncomplete = map[string]Uncomplete{}

// TODO: 优先级队列
// 服务端接收到消息(请求)队列
var recv_queue = []*MessagePack{}

// 服务端返回消息(通知)队列
var send_queue = []*MessagePack{}

func messageSize(in []byte) int {
	return int(binary.BigEndian.Uint16(in))
}

func MessageForward(msg *MessagePack) {
	//fmt.Println(len(msg.Data))
}

func PushMore(conn_id string, in []byte, size int) {
	if size == 1 {
		// 长度不明确的包
		umsg := Uncomplete{
			Message: &MessagePack{
				Id:   conn_id,
				Data: in,
			},
			Read:   0,
			Header: in,
		}

		uncomplete[conn_id] = umsg

	} else {
		// 长度明确的包
		msg_size := messageSize(in)
		if size >= msg_size+2 {
			// 完整的包
			msg := MessagePack{
				Id:   conn_id,
				Size: msg_size,
				Data: in[2 : msg_size+2],
			}

			MessageForward(&msg)

			if size > msg_size+2 {
				// 超出部分处理
				PushMore(conn_id, in[msg_size+2:], size-msg_size-2)
			}

		} else {
			// 不完整的包
			umsg := Uncomplete{
				Message: &MessagePack{
					Id:   conn_id,
					Size: msg_size,
					Data: in[2:],
				},
				Read:   len(in[2:]),
				Header: in[:2],
			}
			uncomplete[conn_id] = umsg
		}
	}
}

func DataFilter(conn_id string, in []byte, size int) {
	umsg, ok := uncomplete[conn_id]
	if ok {
		// 长度不明确的包补全
		if len(umsg.Header) == 1 {
			umsg.Header = append(umsg.Header, in[0])
			umsg.Message.Size = messageSize(umsg.Header)
			in = in[1:]
			size = size - 1
		}

		left := umsg.Message.Size - umsg.Read
		if left <= size {
			umsg.Message.Data = append(umsg.Message.Data, in[:left]...)
			MessageForward(umsg.Message)
			// free
			delete(uncomplete, conn_id)

			// push more
			if size-left > 0 {
				PushMore(conn_id, in[left:], size-left)
			}

		} else {
			umsg.Message.Data = append(umsg.Message.Data, in...)
			umsg.Read += size
		}

	} else {
		PushMore(conn_id, in, size)
	}
}
