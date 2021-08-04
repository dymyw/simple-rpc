package main

import (
    "fmt"
    "github.com/dymyw/simple-rpc/internal/client"
    "github.com/dymyw/simple-rpc/internal/server"
    "net"
    "time"
)

func main() {
    addr := "localhost:9009"
    // 创建服务
    srv := server.NewServer(addr)
    // 服务注册
    srv.Register("QueryUser", QueryUser)
    // 服务启动
    go srv.Run()

    time.Sleep(1 * time.Second)

    // 客户端连接
    address, _ := net.ResolveTCPAddr("tcp", addr)
    conn, err := net.DialTCP("tcp", nil, address)
    defer conn.Close()
    if err != nil {
        panic(err)
    }
    cli := client.NewClient(conn)

    // 桩代码
    var Query func(int) (User, error)
    // 客户端打桩
    cli.CallPRC("QueryUser", &Query)

    // 远程调用
    u, err := Query(5)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Name: %s, Age: %d \n", u.Name, u.Age)
    // Name: 王五, Age: 35
}
