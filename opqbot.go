package main

import (
	"github.com/mcoo/OPQBot"
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
