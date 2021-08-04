package server

// todo NameService 注册服务中心
// todo 异步设计
// todo 双工通信

import (
    "fmt"
    "github.com/dymyw/simple-rpc/internal/dataserial"
    "github.com/dymyw/simple-rpc/internal/transport"
    "io"
    "log"
    "net"
    "reflect"
)

// RPCServer 服务
type RPCServer struct {
    addr string
    funcs map[string]reflect.Value
}

// NewServer 创建服务
func NewServer(addr string) *RPCServer {
    return &RPCServer{
        addr: addr,
        funcs: make(map[string]reflect.Value),
    }
}

// Register 服务注册
func (s *RPCServer) Register(fnName string, fnFunc interface{}) {
    if _, ok := s.funcs[fnName]; ok {
        return
    }

    // 获取方法值
    s.funcs[fnName] = reflect.ValueOf(fnFunc)
}

// Execute 执行方法，如果存在
func (s *RPCServer) Execute(req dataserial.RPCdata) dataserial.RPCdata {
    // 获取方法
    f, ok := s.funcs[req.Name]
    if !ok {
        err := fmt.Sprintf("func %s not Registered", req.Name)
        log.Println(err)

        return dataserial.RPCdata{Name: req.Name, Err: err}
    }

    log.Printf("func %s is called\n", req.Name)

    // 获取参数
    inArgs := make([]reflect.Value, len(req.Args))
    for i := range req.Args {
        inArgs[i] = reflect.ValueOf(req.Args[i])
    }

    // 调用
    out := f.Call(inArgs)

    // 获取结果
    resArgs := make([]interface{}, len(out)-1)
    for i := 0; i < len(out)-1; i++ {
        resArgs[i] = out[i].Interface()
    }

    // 获取错误
    var e string
    if _, ok := out[len(out)-1].Interface().(error); ok {
        // 断言获取错误串
        e = out[len(out)-1].Interface().(error).Error()
    }

    return dataserial.RPCdata{Name: req.Name, Args: resArgs, Err: e}
}

func (s *RPCServer) Run() {
    addr, _ := net.ResolveTCPAddr("tcp", s.addr)
    listener, err := net.ListenTCP("tcp", addr)
    defer listener.Close()
    if err != nil {
        log.Printf("listen on %s err: %v\n", s.addr, err)
        return
    }

    for {
        conn, err := listener.AcceptTCP()
        if err != nil {
            log.Printf("accept err: %v\n", err)
            continue
        }

        go func() {
            connTransport := transport.NewTransport(conn)

            for {
                // 读取请求数据
                req, err := connTransport.Read()
                if err != nil {
                    if err != io.EOF {
                        log.Printf("read err: %v\n", err)
                        return
                    }
                }

                // 反序列化
                decReq, err := dataserial.Decode(req)
                if err != nil {
                    log.Printf("Error Decoding the Payload err: %v\n", err)
                    return
                }

                // 执行
                resP := s.Execute(decReq)

                // 序列化结果
                b, err := dataserial.Encode(resP)
                if err != nil {
                    log.Printf("Error Encoding the Payload for response err: %v\n", err)
                    return
                }

                // 响应
                err = connTransport.Send(b)
                if err != nil {
                    log.Printf("transport write err: %v\n", err)
                }
            }
        }()
    }
}
