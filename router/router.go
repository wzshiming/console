package router

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/wzshiming/console"
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

func newWsConn(conn *websocket.Conn) *wsConn {
	return &wsConn{
		conn: conn,
	}
}

func (c *wsConn) Read(p []byte) (n int, err error) {
	_, rc, err := c.conn.NextReader()
	if err != nil {
		return 0, err
	}
	return rc.Read(p)
}

func (c *wsConn) Write(p []byte) (n int, err error) {
	wc, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}
	defer wc.Close()
	return wc.Write(p)
}

func ResponseJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return err
	}
	return nil
}

func ExecRouter(disable bool, con *console.ReqCreateExec) *mux.Router {
	// 路由
	mux0 := mux.NewRouter()
	upgrader := websocket.Upgrader{
		EnableCompression: true,
	}

	// 创建连接
	mux0.HandleFunc("/create_exec", func(w http.ResponseWriter, r *http.Request) {
		req := &console.ReqCreateExec{}
		*req = *con
		if !disable {
			err := requests(r.Body, &req)
			if err != nil {
				ResponseJSON(w, http.StatusBadRequest, errMsg{err.Error()})
				return
			}
		}

		// 设置默认连接驱动
		if req.Name == "" {
			u, err := url.Parse(req.Host)
			if err != nil {
				ResponseJSON(w, http.StatusBadRequest, errMsg{err.Error()})
				return
			}
			req.Name = u.Scheme
		}

		// 获取驱动
		sesss, err := console.GetDrivers(req.Name, req.Host)
		if err != nil {
			ResponseJSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		// 创建连接
		exec, err := sesss.CreateExec(req)
		if err != nil {
			ResponseJSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		addSession(exec.EId, sesss)

		ResponseJSON(w, http.StatusOK, exec)
		return
	})

	// 开始连接
	mux0.HandleFunc("/start_exec", func(w http.ResponseWriter, r *http.Request) {

		eid := r.FormValue("eid")

		client := getSession(eid)
		if client == nil {
			ResponseJSON(w, http.StatusBadRequest, nil)
			return
		}
		defer delSession(eid)

		if !websocket.IsWebSocketUpgrade(r) {
			ResponseJSON(w, http.StatusBadRequest, nil)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			ResponseJSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}
		defer ws.Close()

		ws.SetCloseHandler(nil)
		ws.SetPingHandler(nil)

		// 执行连接
		err = client.StartExec(eid, newWsConn(ws))
		if err != nil {
			ResponseJSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}
		return
	})

	// 窗口大小调整
	mux0.HandleFunc("/resize_exec_tty", func(w http.ResponseWriter, r *http.Request) {
		req := &console.ReqResizeExecTTY{}
		err := requests(r.Body, &req)
		if err != nil {
			ResponseJSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		client := getSession(req.EId)
		if client == nil {
			ResponseJSON(w, http.StatusBadRequest, nil)
			return
		}

		err = client.ResizeExecTTY(req)
		if err != nil {
			ResponseJSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		ResponseJSON(w, http.StatusOK, nil)
		return
	})

	return mux0
}
