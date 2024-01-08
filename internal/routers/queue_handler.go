package routers

import (
	"net/http"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type Queue struct {
}

func (t *Queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, cancel := results.Get()
	defer cancel()
	nctx := r.Context()

	bt, err := redisx.QueueInfo(nctx, redisx.Client(), redisx.QueueKey(redisx.BqConfig.Redis.Prefix))
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