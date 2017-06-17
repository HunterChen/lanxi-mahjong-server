package socket

import (
	"game/interfacer"
	"net"
	"net/http"
	"runtime/debug"
	"time"
	"net/http/pprof"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"runtime"
)

var (
	VERSION    = "0.0.1"
	BUILD_TIME = ""
	RUN_TIME   = time.Now().Format("2006-01-02 15:04:05")
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

func printroominfo(w http.ResponseWriter, r *http.Request) {

	//body, _ := ioutil.ReadAll(c.Request().Body)
	//roomInfoReq := &cheat.RoomInfoReq{}
	//err := json.Unmarshal(body, roomInfoReq)
	//if err != nil {
	//	return c.JSON(http.StatusOK, H{"code": 1007})
	//}

	//r := room.Get(strconv.Itoa(roomInfoReq.RoomId))
	//
	//if r == nil {
	//	w.Write([]byte("房间不存在"))
	//}
	//
	//w.Write([]byte(r.ToString()))
}

func release(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("go version：" + runtime.Version() +
		"\n build time: " + BUILD_TIME +
		"\n version: " + VERSION +
		"\n startup time: " + RUN_TIME ))
}
func routes() (r *mux.Router) {
	r = mux.NewRouter()
	r.HandleFunc("/", wSHandler).Methods("GET")
	r.HandleFunc("/release", release)
	r.HandleFunc("/printroominfo/", printroominfo)

	// debug
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol).Methods("POST")
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	//r.HandleFunc("/debug/pprof/", http.HandlerFunc(pprof.Index))
	//r.HandleFunc("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	//r.HandleFunc("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	//r.HandleFunc("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	//r.HandleFunc("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

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
		case <-h.closeChan:
			for _, c := range h.connections {
				c.Close()
			}
		}
	}
}
