package v1

import (
	"encoding/json"
	"io"
	"net/http"
)

// Context 如果不需要拓展写成结构体
type Context struct {
	W http.ResponseWriter
	R *http.Request
}

// ReadJson 读取body，反序列化
func (c *Context) ReadJson(req interface{}) error {
	r := c.R.Body
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		return err
	}
	//这里return nil 他可以直接修改地址上的结构体这么理解
	return nil
}
func (c *Context) writeJson(code int, resp interface{}) error {
	c.W.WriteHeader(code)
	result, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.W.Write(result)
	return err
}
func (c *Context) OkJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.writeJson(http.StatusOK, data)
}

func (c *Context) SystemErrJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.writeJson(http.StatusInternalServerError, data)
}

func (c *Context) BadRequestJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.writeJson(http.StatusBadRequest, data)
}
