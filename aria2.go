package main

import (
	"context"
	"github.com/mcoo/OPQBot"
	"github.com/zyxar/argo/rpc"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func pr(a int) string {
	var pr string
	for i := 0; i < a; i++ {
		pr += "█"
	}
	return pr
}

type aria2c struct {
	url, token *string
	a          rpc.Client
}

func (a *aria2c) Connaria2() error {
	ctx := context.Background()
	var no rpc.Notifier
	rsp, err := rpc.New(ctx, *a.url, *a.token, time.Second*20, no)
	if err != nil {
		return err
	}
	a.a = rsp
	ver, err := rsp.GetVersion()
	log.Println("connect successful,ver:", ver.Version)
	return nil
}
func (a *aria2c) ondown(gid string, groupid int64, userid int64, opqbot *OPQBot.BotManager) {
	for true {
		rsp, err := a.a.TellStatus(gid)
		if err != nil {
			log.Println("ondown err:", err)
			break
		}
		if rsp.Status == "complete" {
			if rsp.BitTorrent.Info.Name != "" {
				send2gp(opqbot, groupid, "[ATUSER("+strconv.FormatInt(userid, 10)+")]\n下载任务完成！\n文件名："+rsp.BitTorrent.Info.Name+"\nGid:"+gid+"\n请扫码获取文件[PICFLAG]", "https://z3.ax1x.com/2021/04/29/gkVsZ8.png")
				limit -= 1
			} else {
				if len(rsp.FollowedBy) == 1 {
					go a.ondown(rsp.FollowedBy[0], groupid, userid, opqbot)
					break
				} else {
					send2gp(opqbot, groupid, "[ATUSER("+strconv.FormatInt(userid, 10)+")]\n下载任务完成！\n文件名："+strings.Trim(rsp.Files[0].Path, rsp.Dir)+"\nGid:"+gid+"\n请扫码获取文件[PICFLAG]", "https://z3.ax1x.com/2021/04/29/gkVsZ8.png")
					limit -= 1
				}
			}
			break
		} else if rsp.Status != "active" {
			if rsp.Status == "error" {
				send2g(opqbot, groupid, "[ATUSER("+strconv.FormatInt(userid, 10)+")]\n下载任务失败！\n文件名："+strings.Trim(rsp.Files[0].Path, rsp.Dir)+"\nGid："+gid+"\nErrMsg："+rsp.ErrorMessage)
			}
			break
		}
		time.Sleep(20 * time.Second)
	}
}
func (a *aria2c) Addurl(url1 string) (string, error) {
	url := make([]string, 1)
	url[0] = url1
	gid, err := a.a.AddURI(url)
	if err != nil {
		return gid, err
		log.Println(err)
	}
	return gid, nil
}

func (a *aria2c) Addbt(url string) (string, error) {
	cmd := exec.Command("wget", url, "-O", "./tmp/tmp.torrent")
	cmd.Run()
	gid, err := a.a.AddTorrent("./tmp/tmp.torrent")
	if err != nil {
		return "err", err
	}
	os.Remove("./tmp/tmp.torrent")
	return gid, nil
}

func (a *aria2c) Filestatus(gid string) (string, error) {
	rsp, err := a.a.TellStatus(gid)
	if err != nil {
		return "err", err
		log.Println(err)
	}
	spi, err := strconv.ParseInt(rsp.DownloadSpeed, 10, 64)
	toi, err := strconv.ParseFloat(rsp.TotalLength, 64)
	cpi, err := strconv.ParseFloat(rsp.CompletedLength, 64)
	if rsp.BitTorrent.Info.Name != "" {
		return "文件名：" + rsp.BitTorrent.Info.Name + "\n下载状态：" + rsp.Status + "\n下载速度：" + strconv.FormatInt(spi/1024, 10) + "KB/s\n下载进度：" + strconv.Itoa(int(cpi/toi*100)) + "%\n" + rsp.CompletedLength + "/" + rsp.TotalLength, err
	} else {
		if len(rsp.FollowedBy) == 1 {
			return a.Filestatus(rsp.FollowedBy[0])
		} else {
			return "文件名：" + strings.Trim(rsp.Files[0].Path, rsp.Dir) + "\n下载状态：" + rsp.Status + "\n下载速度：" + strconv.FormatInt(spi/1024, 10) + "KB/s\n下载进度：" + strconv.Itoa(int(cpi/toi*100)) + "%\n" + rsp.CompletedLength + "/" + rsp.TotalLength, err
		}
	}
}

func (a *aria2c) Stop(gid string) error {
	_, err := a.a.Pause(gid)
	return err
}

func (a *aria2c) Start(gid string) error {
	_, err := a.a.Unpause(gid)
	return err
}

func (a *aria2c) Del(gid string) error {
	_, err := a.a.Remove(gid)
	return err
}
func (a *aria2c) Closearia2() {
	a.a.Close()
}
