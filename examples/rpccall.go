package main

// todo test

import (
    "encoding/gob"
    "fmt"
    "github.com/dymyw/simple-rpc/internal/client"
    "github.com/dymyw/simple-rpc/internal/server"
    "net"
    "time"
)

// User1 传输对象
type User1 struct {
    Name string
    Age int
}

// UserDB1 实例
var userDB1 = map[int]User1{
    3: User1{"张三", 23},
    6: User1{"赵六", 66},
    5: User1{"王五", 35},
}

// QueryUser1 执行方法
func QueryUser1(id int) (User1, error) {
    if u, ok := userDB1[id]; ok {
        return u, nil
    }

    return User1{}, fmt.Errorf("id %d not in user db", id)
}

func main() {
    gob.Register(User1{})
    addr := "localhost:9008"
    // 创建服务
    srv := server.NewServer(addr)
    // 服务注册
    srv.Register("QueryUser", QueryUser1)
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
    var Query func(int) (User1, error)
    // 客户端打桩
    cli.CallPRC("QueryUser", &Query)

    // 远程调用
    u, err := Query(6)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Name: %s, Age: %d \n", u.Name, u.Age)
    // Name: 王五, Age: 35
}
