package routers

import (
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/spf13/viper"
)

type Queue struct {
	client *redis.Client
}

func NewQueue(client *redis.Client) *Queue {
	return &Queue{client: client}
}

func (t *Queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, cancel := results.Get()
	defer cancel()

	// url like: queue?list&page=1&pageSize=10
	query := r.URL.RawQuery
	querys := strings.Split(query, "&")
	if len(querys) < 1 {
		result.Code = "1004"
		result.Msg = "404"
		_ = result.Json(w, http.StatusNotFound)
		return
	}

	action := querys[0]
	if r.Method == http.MethodGet {
		// queue list
		if action == "list" {
			bt, err := redisx.QueueInfo(r.Context(), t.client, redisx.QueueKey(redisx.BqConfig.Redis.Prefix))
			if err != nil {
				result.Code = consts.InternalServerErrorCode
				result.Msg = err.Error()
				_ = result.Json(w, http.StatusInternalServerError)
				return
			}

			result.Data = bt
			_ = result.Json(w, http.StatusOK)
			return
		}
		// queue detail
		if action == "detail" {
			queueDetail(w, r, t.client)
		}
	}
}

func queueDetail(w http.ResponseWriter, r *http.Request, client *redis.Client) {

	result, cancel := results.Get()
	defer cancel()

	flusher, ok := w.(http.Flusher)
	if !ok {

	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	id := r.FormValue("id")
	prefix := viper.GetString("redis.prefix")
	id = strings.Join([]string{prefix, id, "stream"}, ":")

	ctx := r.Context()
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			cmd, err := client.XInfoStreamFull(ctx, id, 10).Result()
			if err != nil {
				result.Code = "1004"
				result.Msg = err.Error()
			}

			if err == nil {
				result.Data = cmd.Entries
			}
			_ = result.EventMsg(w, "queue_detail")
			flusher.Flush()
			ticker.Reset(10 * time.Second)
		}
	}
}