package routers

import (
	"net/http"
	"strings"

	"github.com/retail-ai-inc/beanqui/internal/jwtx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

func Auth(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		result, cancel := results.Get()
		defer cancel()

		auth := request.Header.Get("Beanq-Authorization")

		strs := strings.Split(auth, " ")
		if len(strs) < 2 {
			// return data format err
			result.Code = consts.InternalServerErrorCode
			result.Msg = "missing parameter"
			_ = result.Json(writer, http.StatusInternalServerError)
			return
		}

		token, err := jwtx.ParseRsaToken(strs[1])
		if err != nil {
			result.Code = consts.InternalServerErrorCode
			result.Msg = err.Error()
			_ = result.Json(writer, http.StatusUnauthorized)
			return
		}
		//
		_, err = token.Claims.GetExpirationTime()
		if err != nil {
			result.Code = consts.InternalServerErrorCode
			result.Msg = err.Error()
			_ = result.Json(writer, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(writer, request)
	}
}
