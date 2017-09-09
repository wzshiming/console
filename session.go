package console

import (
	"io"
)

type Sessions interface {
	CreateExec(req *ReqCreateExec) (*RespCreateExec, error)
	StartExec(id string, ws io.ReadWriter) error
	ResizeExecTTY(req *ReqResizeExecTTY) error
}

type ReqCreateExec struct {
	Name string `json:"name,omitempty"`
	Host string `json:"host,omitempty"`
	CId  string `json:"cid,omitempty"`
	Cmd  string `json:"cmd,omitempty"`
}

type RespCreateExec struct {
	EId string `json:"eid,omitempty"`
}

type ReqResizeExecTTY struct {
	EId    string `json:"eid,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}
