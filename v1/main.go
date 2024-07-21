package v1

import (
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "这是主页")
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "这是用户")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "这是创建用户")
}

func order(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "这是订单")
}
func main() {
	//优化1 这里会有很多路由和对应handler的操作,以及启动，因此这里可以抽象出一个接口

	server := NewServer("test-server")
	//server.Route("/", home)
	//server.Route("/user", user)
	//server.Route("/user/create", createUser)

	//增加RestFul

	server.Route("POST", "/user/signup", SignUp)
	//server.Route("/order", order)
	err := server.Start(":8080")
	if err != nil {
		return
	}
}
