package notify

import (
	"bytes"
	"log"
	"net/http"
	"sync"
	"time"
)

// 前台 Next.js 客户端的 ISR 缓存失效通知。
// 未配置 revalidate_url 时所有调用都是空操作，不产生任何网络请求。
var (
	mu              sync.Mutex
	revalidateURL   string
	revalidateToken string
	inFlight        bool
	dirty           bool
	client          = &http.Client{Timeout: 5 * time.Second}
)

func Init(url, token string) {
	mu.Lock()
	defer mu.Unlock()
	revalidateURL = url
	revalidateToken = token
}

// RevalidateClient 异步通知前台全量失效缓存；短时间内的多次调用会合并为一次请求。
func RevalidateClient() {
	mu.Lock()
	if revalidateURL == "" {
		mu.Unlock()
		return
	}
	dirty = true
	if inFlight {
		mu.Unlock()
		return
	}
	inFlight = true
	mu.Unlock()

	go func() {
		for {
			mu.Lock()
			if !dirty {
				inFlight = false
				mu.Unlock()
				return
			}
			dirty = false
			url, token := revalidateURL, revalidateToken
			mu.Unlock()
			send(url, token)
		}
	}()
}

func send(url, token string) bool {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(`{"type":"all"}`)))
	if err != nil {
		log.Printf("[notify] build revalidate request: %v", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-revalidate-secret", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[notify] revalidate client: %v", err)
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

// RevalidateOnStartup 在后端启动后多次尝试通知前台失效缓存。
// 前台镜像构建时通常连不上后端，预渲染页面是空数据；
// 部署完成后这一次通知能让前台立即用真实数据重建缓存。
func RevalidateOnStartup() {
	mu.Lock()
	url, token := revalidateURL, revalidateToken
	mu.Unlock()
	if url == "" {
		return
	}
	for i := 0; i < 6; i++ {
		time.Sleep(5 * time.Second)
		if send(url, token) {
			return
		}
	}
}
