package main

import "fmt"

// User 传输对象
type User struct {
    Name string
    Age int
}

// UserDB 实例
var userDB = map[int]User{
    3: User{"张三", 23},
    6: User{"赵六", 66},
    5: User{"王五", 35},
}

// QueryUser 执行方法
func QueryUser(id int) (User, error) {
    if u, ok := userDB[id]; ok {
        return u, nil
    }

    return User{}, fmt.Errorf("id %d not in user db", id)
}

func main() {
    u, err := QueryUser(5)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("Name: %s, Age: %d \n", u.Name, u.Age)
    // Name: 王五, Age: 35
}
