package client

// 客户端例子
//  var rpcTest func(id int) (int, error)
//  MakeRpc("testrpc", &rpcTest)
//  id, err := rpcTest(10)

import (
    "net"
    "reflect"
)

// Client 客户端
type Client struct {
    conn *net.TCPConn
}

// NewClient 实例化
func NewClient(conn *net.TCPConn) *Client {
    return &Client{conn}
}

// CallPRC 远程调用
func (c *Client) CallPRC(rpcName string, fPtr interface{}) {
    // 取值指向的元素值
    container := reflect.ValueOf(fPtr).Elem()

    // 桩代码
    f := func(req []reflect.Value) []reflect.Value {
        // todo

        return []reflect.Value{}
    }

    // 设置反射对象值
    container.Set(reflect.MakeFunc(container.Type(), f))
}
