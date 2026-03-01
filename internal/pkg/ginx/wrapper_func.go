package ginx

import "github.com/gin-gonic/gin"

func WrapBodyAndToekn[Req any, C any](bizFn func(ctx *gin.Context, req Req, uc C)) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
