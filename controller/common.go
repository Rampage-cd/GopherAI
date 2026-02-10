package controller

import "GopherAI/common/code"

type Response struct {
	StatusCode code.Code `json:"status_code"`
	StatusMsg  string    `json:"status_msg,omitempty"`
} //用来记录状态码，和状态信息

func (r *Response) CodeOf(code code.Code) *Response {
	if r == nil {
		r = &Response{}
	}
	r.StatusCode = code
	r.StatusMsg = code.Msg() //状态码对应的状态信息
	return r
} //创建结构体并存放状态码和状态信息

func (r *Response) Success() {
	r.CodeOf(code.CodeSuccess)
} //存放成功的状态信息
