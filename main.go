package main

import (
	"encoding/json"
	"fmt"
	"github.com/mcoo/OPQBot"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var limit int //下载任务限制
var ver = "Aria2M_for_OPQ_ver.0.2b"

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func main() {
	limit = 0
	fmt.Println(ver)
	fmt.Println("By Liumik")
	if !Exists("./config.json") {
		tmp := make(map[string]interface{})
		var url string
		var token string
		var qq int64
		var site string
		fmt.Println("\n请输入OPQ的Web地址: ")
		fmt.Scan(&site)
		fmt.Println("\n请输入Bot账号: ")
		fmt.Scan(&qq)
		fmt.Println("\n请输入url")
		fmt.Scan(&url)
		fmt.Println("请输入token")
		fmt.Scan(&token)
		tmp["Site"] = site
		tmp["Qq"] = qq
		tmp["Url"] = url
		tmp["Token"] = token
		tmp1, _ := json.Marshal(tmp)
		c1, err := os.Create("config.json")
		defer c1.Close()
		if err != nil {
			log.Println("cerr:", err)
			os.Exit(1)
		}
		c1.Write(tmp1)
	}
	c1, err := os.OpenFile("./config.json", os.O_RDONLY, 0600)
	defer c1.Close()
	if err != nil {
		log.Println("openerr:", err)
		os.Exit(1)
	}
	cb, _ := ioutil.ReadAll(c1)
	conf1 := struct {
		Site  string `Site`
		Qq    int64  `Qq`
		Url   string `Url`
		Token string `Token`
	}{}
	json.Unmarshal(cb, &conf1)
	opqBot := OPQBot.NewBotManager(conf1.Qq, conf1.Site)
	err1 := opqBot.Start()
	if err1 != nil {
		log.Println(err.Error())
	}
	defer opqBot.Stop()
	ac1 := aria2c{token: &conf1.Token, url: &conf1.Url}
	err2 := ac1.Connaria2()
	if err2 != nil {
		log.Fatalln("aria2 error:", err2)
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet OPQBot.GroupMsgPack) {
		//log.Println(botQQ, packet.Content)
		fileinfo := struct {
			FileID   string `FileID`
			FileName string `FileName`
		}{}
		json.Unmarshal([]byte(packet.Content), &fileinfo)
		if strings.HasPrefix(packet.Content, "addurl_") {
			if limit <= 5 {
				gid, err := ac1.Addurl(strings.TrimPrefix(packet.Content, "addurl_"))
				if err != nil {
					send2g(&opqBot, packet.FromGroupID, "error:"+err.Error())
				} else {
					send2g(&opqBot, packet.FromGroupID, "已添加下载任务，发送status_"+gid+"查看详情")
					limit += 1
					go ac1.ondown(gid, packet.FromGroupID, packet.FromUserID, &opqBot)
				}
			} else {
				send2g(&opqBot, packet.FromGroupID, "下载任务过多，请等待任务完成")
			}
		}
		if strings.HasPrefix(packet.Content, "status_") {
			rsp, err := ac1.Filestatus(strings.TrimPrefix(packet.Content, "status_"))
			if err != nil {
				send2g(&opqBot, packet.FromGroupID, "error:"+err.Error())
			} else {
				send2g(&opqBot, packet.FromGroupID, rsp)
			}
		}
		if strings.HasPrefix(fileinfo.FileName, "addbt_") {
			if limit <= 5 {
				_, urlt, err := opqBot.GetFile(fileinfo.FileID, packet.FromGroupID)
				if err != nil {
					log.Println(err)
				}
				gid, err := ac1.Addbt(urlt.URL)
				if err != nil {
					send2g(&opqBot, packet.FromGroupID, "error:"+err.Error())
				} else {
					send2g(&opqBot, packet.FromGroupID, "已添加下载任务，发送status_"+gid+"查看详情")
					limit += 1
					go ac1.ondown(gid, packet.FromGroupID, packet.FromUserID, &opqBot)
				}
			} else {
				send2g(&opqBot, packet.FromGroupID, "下载任务过多，请等待任务完成")
			}
		}
		if strings.HasPrefix(packet.Content, "stop_") {
			err := ac1.Stop(strings.TrimPrefix(packet.Content, "stop_"))
			if err != nil {
				send2g(&opqBot, packet.FromGroupID, "error:"+err.Error())
			} else {
				send2g(&opqBot, packet.FromGroupID, "已暂停下载任务")
			}
		}
		if strings.HasPrefix(packet.Content, "start_") {
			err := ac1.Start(strings.TrimPrefix(packet.Content, "start_"))
			if err != nil {
				send2g(&opqBot, packet.FromGroupID, "error:"+err.Error())
			} else {
				send2g(&opqBot, packet.FromGroupID, "已开始下载任务")
			}

		}
		if strings.HasPrefix(packet.Content, "del_") {
			err := ac1.Del(strings.TrimPrefix(packet.Content, "del_"))
			if err != nil {
				send2g(&opqBot, packet.FromGroupID, "error:"+err.Error())
			} else {
				send2g(&opqBot, packet.FromGroupID, "已移除下载任务")
			}
		}
	})
	err = opqBot.AddEvent(OPQBot.EventNameOnFriendMessage, func(botQQ int64, packet *OPQBot.FriendMsgPack) {
		log.Println(botQQ, packet.Content)
	})

	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnConnected, func() {
		log.Println("连接成功！！！")
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnDisconnected, func() {
		log.Println("连接断开！！")
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = opqBot.AddEvent(OPQBot.EventNameOnOther, func(botQQ int64, e interface{}) {
		log.Println(e)
	})
	if err != nil {
		log.Println(err.Error())
	}
	opqBot.Wait()
}
