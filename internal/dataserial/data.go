package dataserial

// todo interface
// todo 常规结构序列化实现
// todo 二进制序列化实现
// todo 专用序列化实现

import (
	"bytes"
	"encoding/gob"
)

// RPCdata 远程调用传输结构
type RPCdata struct {
	Name string			// 方法名
	Args []interface{}	// 参数
	Err string			// 运行时错误
}

// Encode 序列化
func Encode(data RPCdata) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decode 反序列化
func Decode(b []byte) (RPCdata, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)

	var data RPCdata
	if err := decoder.Decode(&data); err != nil {
		return RPCdata{}, err
	}

	return data, nil
}
