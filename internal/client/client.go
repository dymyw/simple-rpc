package client

// 客户端例子
//  var rpcTest func(id int) (int, error)
//  MakeRpc("testrpc", &rpcTest)
//  id, err := rpcTest(10)

import (
    "errors"
    "github.com/dymyw/simple-rpc/internal/dataserial"
    "github.com/dymyw/simple-rpc/internal/transport"
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
        // 请求错误处理
        errorHandler := func(err error) []reflect.Value {
            // 出参
            numOut := container.Type().NumOut()
            outArgs := make([]reflect.Value, numOut)

            for i := 0; i < len(outArgs)-1; i++ {
                outArgs[i] = reflect.Zero(container.Type().Out(i))
            }
            outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()

            return outArgs
        }

        // 参数
        inArgs := make([]interface{}, 0, len(req))
        for _, arg := range req {
            // 转换
            inArgs = append(inArgs, arg.Interface())
        }

        // 请求
        reqData := dataserial.RPCdata{Name: rpcName, Args: inArgs}
        b, err := dataserial.Encode(reqData)
        if err != nil {
            panic(err)
        }
        cReqTransport := transport.NewTransport(c.conn)
        err = cReqTransport.Send(b)
        if err != nil {
            return errorHandler(err)
        }

        // 读取
        rsp, err := cReqTransport.Read()
        if err != nil {
            return errorHandler(err)
        }
        rspData, _ := dataserial.Decode(rsp)
        // 错误处理
        if rspData.Err != "" {
            return errorHandler(errors.New(rspData.Err))
        }
        // 无返回值
        if len(rspData.Args) == 0 {
            rspData.Args = make([]interface{}, container.Type().NumOut())
        }

        // 出参
        numOut := container.Type().NumOut()
        outArgs := make([]reflect.Value, numOut)
        for i := 0; i < numOut; i++ {
            // 转换
            if i != numOut-1 {
                if rspData.Args[i] == nil {
                    outArgs[i] = reflect.Zero(container.Type().Out(i))
                } else {
                    outArgs[i] = reflect.ValueOf(rspData.Args[i])
                }
            } else {
                // Err
                outArgs[i] = reflect.Zero(container.Type().Out(i))
            }
        }

        return outArgs
    }

    // 设置反射对象值
    container.Set(reflect.MakeFunc(container.Type(), f))
}
