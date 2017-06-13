/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 10:23
 * Filename      : hub.go
 * Description   : 负责用户登入登出和用户间消息广播
 * *******************************************************/
package socket

import (
	"game/interfacer"
	"net"
	"net/http"
	"runtime/debug"
	//"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

type broadcastPacket struct {
	userids     []string
	content     interfacer.IProto
	successChan chan []string
}
type detectOnline struct {
	userids    []string
	detectChan chan []interfacer.IConn
}

func routes() (r *mux.Router) {
	r = mux.NewRouter()
	r.HandleFunc("/", wSHandler).Methods("GET")
	return
}

func Server(addr string) (ln net.Listener, ch chan error) {
	go h.run()
	ch = make(chan error)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		r := routes()
		ch <- http.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)}, r)
	}()
	return
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func logout(c interfacer.IConn) {
	h.unregister <- c
	// sync online table
}

// 在线人数
func OnlineCount() uint32 {
	callback := make(chan uint32)
	h.onlineCount <- callback
	return <-callback
}

// 在线人数
func Close() {
	h.closeChan <- true
}

type hub struct {
	connections      map[string]interfacer.IConn
	broadcast        chan *broadcastPacket
	register         chan interfacer.IConn
	unregister       chan interfacer.IConn
	detectonlineChan chan *detectOnline
	onlineCount      chan chan uint32
	closeChan        chan bool
}

var h = hub{
	connections:      make(map[string]interfacer.IConn, 1024),
	broadcast:        make(chan *broadcastPacket, 1024),
	register:         make(chan interfacer.IConn, 32),
	unregister:       make(chan interfacer.IConn, 32),
	detectonlineChan: make(chan *detectOnline, 32),
	onlineCount:      make(chan chan uint32, 32),
	closeChan:        make(chan bool, 1),
}

func (h *hub) run() {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(string(debug.Stack()))
		}
	}()

	for {
		select {
		case n := <-h.onlineCount:
			n <- uint32(len(h.connections))
		case d := <-h.detectonlineChan:
			users := make([]interfacer.IConn, 0, len(d.userids))
			for _, v := range d.userids {
				con, ok := h.connections[v]
				if ok {
					users = append(users, con)
				}
			}
			d.detectChan <- users
		case c := <-h.register:
			h.connections[c.GetUserid()] = c
		case c := <-h.unregister:
			if conn, ok := h.connections[c.GetUserid()]; ok {
				if conn == c {
					delete(h.connections, c.GetUserid())
				}
			}
		case m := <-h.broadcast:
			if m != nil {
				result := make([]string, 0, len(m.userids))
				for _, v := range m.userids {
					if con, ok := h.connections[v]; ok {
						con.Send(m.content)
						glog.Infoln(m.content)
					} else {
						result = append(result, v)
					}
				}
				m.successChan <- result
			}
		case c := <-h.closeChan:
			//TODO:退出处理
			glog.Infoln("server closed -> ", c, " conns -> ", len(h.connections))
			for _, c := range h.connections {
				c.Close()
			}
		}
	}
}
