package main

import (
	"bytes"
	"encoding/json"
	"github.com/mcoo/OPQBot"
	"io/ioutil"
	"log"
	"net/http"
)

func send2f(m *OPQBot.BotManager, uid int64, content string) {
	m.Send(OPQBot.SendMsgPack{
		ToUserUid:  uid,
		Content:    OPQBot.SendTypeTextMsgContent{Content: content},
		SendToType: OPQBot.SendToTypeFriend,
	}) //发送好友消息
}
func send2p(m *OPQBot.BotManager, uid int64, gid int64, content string) {
	m.Send(OPQBot.SendMsgPack{
		ToUserUid:  uid,
		Content:    OPQBot.SendTypeTextMsgContentPrivateChat{Content: content, Group: gid},
		SendToType: OPQBot.SendToTypePrivateChat,
	}) //发送群私聊消息
}
func send2g(m *OPQBot.BotManager, uid int64, content string) {
	m.Send(OPQBot.SendMsgPack{
		ToUserUid:  uid,
		Content:    OPQBot.SendTypeTextMsgContent{Content: content},
		SendToType: OPQBot.SendToTypeGroup,
	}) //发送群消息
}
func send2gp(m *OPQBot.BotManager, uid int64, content string, picurl string) {
	m.Send(OPQBot.SendMsgPack{
		ToUserUid:  uid,
		Content:    OPQBot.SendTypePicMsgByUrlContent{Content: content, PicUrl: picurl},
		SendToType: OPQBot.SendToTypeGroup,
	}) //发送群消息
}
func Getfile(groupid int64, fileid string, qq string, url1 string) string {
	url := struct {
		Url string `Url`
	}{}
	tmp := make(map[string]interface{})
	tmp["GroupID"] = groupid
	tmp["FileID"] = fileid
	tmp1, _ := json.Marshal(tmp)
	resp, err := (http.Post(url1+"/v1/LuaApiCaller?funcname=OidbSvc.0x6d6_2&timeout=10&qq="+qq, "application/json", bytes.NewBuffer(tmp1)))
	if err != nil {
		log.Fatal(err)
		return "err"
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	json.Unmarshal(body, &url)
	return url.Url
}
