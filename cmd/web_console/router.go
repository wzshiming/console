package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/unrolled/render"
	"github.com/wzshiming/console"
	_ "github.com/wzshiming/console/docker"
	_ "github.com/wzshiming/console/shell"
	_ "github.com/wzshiming/console/ssh"
)

var (
	sessionsMu sync.RWMutex
	session    = map[string]console.Sessions{}
)

func addSession(id string, s console.Sessions) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	session[id] = s
}

func getSession(id string) console.Sessions {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()
	return session[id]
}

func delSession(id string) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	delete(session, id)
}

func requests(rc io.ReadCloser, i interface{}) error {
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}
	defer rc.Close()
	return json.Unmarshal(data, i)
}

type errMsg struct {
	Msg string `json:"msg,omitempty"`
}

type wsConn struct {
	conn *websocket.Conn
}

func (c *wsConn) Read(p []byte) (n int, err error) {
	_, rc, err := c.conn.NextReader()
	if err != nil {
		return 0, err
	}
	return rc.Read(p)
}

func (c *wsConn) Write(p []byte) (n int, err error) {
	wc, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	defer wc.Close()
	return wc.Write(p)
}

func execRouter() (*mux.Router, error) {
	// 路由
	mux0 := mux.NewRouter()
	rend := render.New()
	upgrader := websocket.Upgrader{
		EnableCompression: true,
	}

	// 创建连接
	mux0.HandleFunc("/create_exec", func(w http.ResponseWriter, r *http.Request) {
		req := &console.ReqCreateExec{}
		err := requests(r.Body, &req)
		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		// 设置默认连接驱动
		if req.Name == "" {
			u, err := url.Parse(req.Host)
			if err != nil {
				rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
				return
			}
			req.Name = u.Scheme
		}

		// 获取驱动
		sesss, err := console.GetDrivers(req.Name, req.Host)
		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		// 创建连接
		exec, err := sesss.CreateExec(req)
		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		addSession(exec.EId, sesss)

		rend.JSON(w, http.StatusOK, exec)
		return
	})

	// 开始连接
	mux0.HandleFunc("/start_exec", func(w http.ResponseWriter, r *http.Request) {

		eid := r.FormValue("eid")

		client := getSession(eid)
		if client == nil {
			rend.JSON(w, http.StatusBadRequest, nil)
			return
		}
		defer delSession(eid)

		if !websocket.IsWebSocketUpgrade(r) {
			rend.JSON(w, http.StatusBadRequest, nil)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}
		defer ws.Close()

		ws.SetCloseHandler(nil)
		ws.SetPingHandler(nil)

		// 执行连接
		err = client.StartExec(eid, &wsConn{ws})
		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		rend.JSON(w, http.StatusSwitchingProtocols, nil)
		return
	})

	// 窗口大小调整
	mux0.HandleFunc("/resize_exec_tty", func(w http.ResponseWriter, r *http.Request) {
		req := &console.ReqResizeExecTTY{}
		err := requests(r.Body, &req)
		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		client := getSession(req.EId)
		if client == nil {
			rend.JSON(w, http.StatusBadRequest, nil)
			return
		}

		err = client.ResizeExecTTY(req)
		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		rend.JSON(w, http.StatusOK, nil)
		return
	})

	return mux0, nil
}
