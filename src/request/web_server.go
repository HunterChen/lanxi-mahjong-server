//web server
package request

import (
	"lib/utils"
	"data"
	"fmt"
	"encoding/json"
	"net/http"
	"players"
	"strings"
	"resource"

	"github.com/gorilla/websocket"
	"github.com/golang/glog"
	"config"
)

type respErr struct {
	ErrCode int    `json:errcode` //错误码
	ErrMsg  string `json:errmsg`  //错误信息
}

type reqMsg struct {
	Userid string `json:userid` //角色ID
	Itemid uint32 `json:itemid` //物品id或类型
	Amount int32  `json:amount` //数量
}

var upgrader = websocket.Upgrader{} // use default options

var WebKey string = "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"

// 支付回调监听服务
func WebServer() {
	if config.Opts().Web_port == "" {
		return
	}

	http.HandleFunc("/", web)
	glog.Fatal(http.ListenAndServe(config.Opts().Web_port, nil))
}

func web(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	if !verifyToken(r.Header.Get("Token")) {
		return
	}
	//
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Infoln("upgrade:", err)
		return
	}
	go webHandler(c) //goroutine
}

func webHandler(c *websocket.Conn) {
	resp := respErr{ErrCode: 0}
	mt, message, err := c.ReadMessage()
	if err != nil {
		resp.ErrCode = 1
		resp.ErrMsg = fmt.Sprintf("read msg err: %v", err)
	}
	ms := strings.Split(string(message), "|")
	var msgJson string
	if len(ms) < 2 || ms[0] != WebKey {
		resp.ErrCode = 2
		resp.ErrMsg = fmt.Sprintf("msg err: %v", err)
	} else {
		msgJson = ms[1]
	}
	glog.Infof("recv: %s", msgJson)
	var msg reqMsg
	err = json.Unmarshal([]byte(msgJson), &msg)
	if err != nil {
		resp.ErrCode = 3
		resp.ErrMsg = fmt.Sprintf("unmarshal msg err: %v", err)
	}
	err = webHandling(msg)
	if err != nil {
		resp.ErrCode = 4
		resp.ErrMsg = fmt.Sprintf("handler msg err: %v", err)
	}
	data, _ := json.Marshal(resp)
	err = c.WriteMessage(mt, data)
	if err != nil {
		glog.Infoln("write:", err)
	}
	c.Close()
}

func webHandling(msg reqMsg) error {
	player := players.Get(msg.Userid)
	if player != nil {
		return resource.ChangeRes(player, msg.Itemid, msg.Amount, data.RESTYPE17)
	}
	user := &data.User{Userid: msg.Userid}
	err := user.Get()
	if err != nil {
		return err
	}
	if msg.Itemid == resource.ROOM_CARD {
		val := map[string]int32{}
		val["RoomCard"] = int32(user.RoomCard) + msg.Amount
		return user.UpdateResource(val)
	}
	return err
}

func verifyToken(Token string) bool {
	// client
	// Key := "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"
	// Now := strconv.FormatInt(utils.Timestamp(), 10)
	// Sign := utils.Md5(Key+Now)
	// Token := Sign+Now+RandNum
	// r.Header.Set("Token")
	// server
	// Key := "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"
	// Token := r.Header.Get("Token")
	// r.Header.Del("Token")
	TokenB := []byte(Token)
	if len(TokenB) >= 42 {
		SignB := TokenB[:32]
		TimeB := TokenB[32:42]
		if utils.Md5(WebKey+string(TimeB)) == string(SignB) {
			return true
		}
	}
	return false
}
