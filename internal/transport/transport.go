package transport

// todo 数据结构化
// todo 请求序号

import (
	"encoding/binary"
	"io"
	"net"
)

// Transport 传输结构体
type Transport struct {
	conn *net.TCPConn
}

// NewTransport 实例化
func NewTransport(conn *net.TCPConn) *Transport {
	return &Transport{conn}
}

// Send 发送数据
func (t *Transport) Send(data []byte) error {
	// 使用 TLV（Type-Length-Value） 协议
	// 4 个字节表示长度，加上数据
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)

	if _, err := t.conn.Write(buf); err != nil {
		return err
	}

	return nil
}

// Read 读取数据
func (t *Transport) Read() ([]byte, error) {
	// 读取长度
	buf := make([]byte, 4)
	if _, err := io.ReadFull(t.conn, buf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(buf)

	// 读取数据
	data := make([]byte, length)
	if _, err := io.ReadFull(t.conn, data); err != nil {
		return nil, err
	}

	return data, nil
}
