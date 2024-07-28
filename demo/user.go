package demo

import (
	"fmt"
	v1 "test/pkg"
	"time"
)

func SignUp(c *v1.Context) {
	req := &signUpReq{}
	err := c.ReadJson(req)
	if err != nil {
		_ = c.BadRequestJson(&commonResponse{
			BizCode: 4, // 这个代表输入参数错误
			Msg:     fmt.Sprintf("invalid request: %v", err),
		})
		return
	}
	_ = c.OkJson(&commonResponse{
		// 这个是新用户的 ID
		Data: 123,
	})
}

func SlowService(c *v1.Context) {
	time.Sleep(time.Second * 10)
	_ = c.OkJson(&commonResponse{
		Msg: "Hi, this is msg from slow service",
	})
}

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
