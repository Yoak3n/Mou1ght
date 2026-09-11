package util

import (
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// 内存版固定窗口限流/去重，按 key 记录窗口内的放行次数。
// 用于点赞、浏览等公开接口的防刷；进程重启即重置，个人博客场景足够。
var (
	rateMu    sync.Mutex
	rateCache = make(map[string]*rateBucket)
)

type rateBucket struct {
	count int
	until time.Time
}

// Allow 在 window 时间窗内最多放行 max 次；超限返回 false。
func Allow(key string, max int, window time.Duration) bool {
	now := time.Now()
	rateMu.Lock()
	defer rateMu.Unlock()
	// 缓存超过阈值时顺带清理过期项，避免无界增长
	if len(rateCache) > 10000 {
		for k, b := range rateCache {
			if now.After(b.until) {
				delete(rateCache, k)
			}
		}
	}
	b := rateCache[key]
	if b == nil || now.After(b.until) {
		rateCache[key] = &rateBucket{count: 1, until: now.Add(window)}
		return true
	}
	if b.count >= max {
		return false
	}
	b.count++
	return true
}

// ClientIP 取真实客户端 IP：优先 X-Forwarded-For 第一个地址（nginx 反代场景），
// 否则使用请求方地址。
func ClientIP(c *fiber.Ctx) string {
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	return c.IP()
}

// VisitorIdentity 解析防刷身份：
//   - 带游客 JWT（Authorization: Bearer ...）时返回其 jti（后端签发、不可伪造）；
//   - 无 token 但来自前台代理（x-proxied-by: client，即浏览器经 Next server action 转发）
//     时返回 tokenRequired=true，表示必须带 token，否则所有访问者会共享同一 IP；
//   - 其余情况（直连 API）回退到真实客户端 IP。
func VisitorIdentity(c *fiber.Ctx) (identity string, isToken bool, tokenRequired bool) {
	auth := c.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		if id, err := ParseVisitorTokenID(strings.TrimPrefix(auth, "Bearer ")); err == nil && id != "" {
			return id, true, false
		}
	}
	if c.Get("x-proxied-by") == "client" {
		return "", false, true
	}
	return ClientIP(c), false, false
}
