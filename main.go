package main

import (
	"strings"
	"time"
	"xyhelper-account-check/config"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

func main() {
	ctx := gctx.New()
	g.Log().Info(ctx, config.CHATPROXY)
	accountSet := gset.New(true)
	gfile.ReadLines("account.txt", func(line string) error {
		// g.Log().Info(ctx, line)
		// 使用,分割
		account := strings.Split(line, config.Separator)
		if len(account) < 2 {
			g.Log().Error(ctx, "账号密码格式错误", line)
			return nil
		}
		accountStr := account[0] + config.Separator + account[1]
		// g.Log().Info(ctx, accountStr)
		accountSet.Add(accountStr)
		return nil
	})
	outputSet := gset.New(true)
	gfile.ReadLines("output.txt", func(line string) error {
		account := strings.Split(line, config.Separator)
		if len(account) < 2 {
			g.Log().Error(ctx, "账号密码格式错误", line)
			return nil
		}
		accountStr := account[0] + config.Separator + account[1]
		// g.Log().Info(ctx, accountStr)
		outputSet.Add(accountStr)
		return nil
	})
	// 求差集 属于accountSet但不属于outputSet
	TaskSet := accountSet.Diff(outputSet)
	g.Log().Info(ctx, "剩余任务数:", TaskSet.Size())
	TaskSet.Walk(func(item interface{}) interface{} {
		// 循环检测status
		for {
			Status := config.StatusCache.MustGet(ctx, "Status").String()
			if Status == "ok" {
				break
			}
			time.Sleep(1 * time.Second)
		}
		// 当符合限速条件时，将任务加入队列
		config.RateLimiter.Wait(ctx)
		g.Log().Info(ctx, item)
		go GetAccountInfo(ctx, item.(string))
		return nil
	})
	g.Log().Info(ctx, "任务已全部加入队列")
	// 阻塞主线程 保持程序运行
	select {}
}

// GetAccountInfo 获取账号信息
func GetAccountInfo(ctx g.Ctx, account string) {
	// 使用,分割 account
	accountInfo := gstr.Split(account, config.Separator)
	if len(accountInfo) < 2 {
		g.Log().Error(ctx, "账号密码格式错误", account)
		return
	}
	username := accountInfo[0]
	password := accountInfo[1]
	resVar := g.Client().PostVar(ctx, config.CHATPROXY+"/applelogin", g.Map{
		"username": username,
		"password": password,
	})
	resJson := gjson.New(resVar)
	// resJson.Dump()
	detail := resJson.Get("detail").String()
	if detail != "" {
		g.Log().Warning(ctx, account, detail)
		result := username + config.Separator + password + config.Separator + detail
		gfile.PutContentsAppend("output.txt", result+"\n")

		return
	}
	refresh_token := resJson.Get("refresh_token").String()
	access_token := resJson.Get("access_token").String()
	if refresh_token == "" || access_token == "" {
		g.Log().Warning(ctx, account, "获取token失败", resVar)
		return
	}
	result := username + config.Separator + password + config.Separator + refresh_token + config.Separator + access_token
	g.Log().Info(ctx, result)
	gfile.PutContentsAppend("output.txt", result+"\n")

}
