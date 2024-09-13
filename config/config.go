package config

import (
	"time"

	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"golang.org/x/time/rate"
)

var (
	CHATPROXY   = ""
	TaskQueue   = gqueue.New()
	TaskSet     = gset.New(true)
	LIMIT       = 1
	PER         = time.Second
	RateLimiter *rate.Limiter
	StatusCache = gcache.New()
)

func init() {
	ctx := gctx.GetInitCtx()
	CHATPROXY = g.Cfg().MustGetWithEnv(ctx, "CHATPROXY").String()
	if CHATPROXY == "" {
		panic("CHATPROXY is empty")
	}
	g.Log().Info(ctx, "CHATPROXY:", CHATPROXY)
	limit := g.Cfg().MustGetWithEnv(ctx, "LIMIT").Int()
	if limit > 0 {
		LIMIT = limit
	}
	per := g.Cfg().MustGetWithEnv(ctx, "PER").Duration()
	g.Dump(per)
	if per > 0 {
		PER = per
	}
	g.Log().Info(ctx, "LIMIT:", LIMIT)
	g.Log().Info(ctx, "PER:", PER)
	RateLimiter = rate.NewLimiter(rate.Every(PER), LIMIT)
	StatusCache.Set(ctx, "LoginSuccessRate", 0, 0)
	CheckStatus(ctx)
	gcron.AddSingleton(ctx, "0 * * * * *", CheckStatus)
}

// CheckStatus 检查状态
func CheckStatus(ctx g.Ctx) {
	// g.Log().Info(ctx, "CheckStatus")
	resVar := g.Client().GetVar(ctx, CHATPROXY+"/ping")
	LoginSuccessRate := gjson.New(resVar).Get("LoginSuccessRate").Int()
	StatusCache.Set(ctx, "LoginSuccessRate", LoginSuccessRate, 0)
	g.Log().Info(ctx, "LoginSuccessRate:", LoginSuccessRate)
	ChatAvailable := gjson.New(resVar).Get("ChatAvailable").Int()
	StatusCache.Set(ctx, "ChatAvailable", ChatAvailable, 0)
	g.Log().Info(ctx, "ChatAvailable:", ChatAvailable)
	if LoginSuccessRate < 80 || ChatAvailable < 1 {
		g.Log().Error(ctx, "LoginSuccessRate <80 || ChatAvailable == 0")
		StatusCache.Set(ctx, "Status", "error", 0)
	} else {
		StatusCache.Set(ctx, "Status", "ok", 0)
	}
}
