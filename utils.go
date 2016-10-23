package gimo

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
)

// dialDB establishes a new session to the mongoDB cluster
// identified by info.
func dialDB(info *mgo.DialInfo) *mgo.Session {
	s, err := mgo.DialWithInfo(info)
	if err != nil {
		panic("Failed to establish session with database: " + err.Error())
	}
	return s
}

// getCtxErrorCode returns the status code that is
// associated with the Gin context error.
func getCtxErrorCode(err *gin.Error) int {
	code := http.StatusInternalServerError
	meta, ok := err.Meta.(*ErrorMeta)
	if ok {
		code = meta.Code
	}
	return code
}

// ErrorsExist returns true is there are errors
// saved on the the Gin context
func ErrorsExist(ctx *gin.Context) bool {
	return len(ctx.Errors) > 0
}

// AbortWithError attaches a Gin error to the current Gin
// context with a status code. It is the same as Gin's
// ctx.AbortWithError function, except that it attaches the
// status code to the context (accompanying the error) instead
// of writing it directly to the HTTP header. This allows a
// separate middleware to handle all the HTTP response writing,
// for the body as well as the header.
func AbortWithError(ctx *gin.Context, err error, code int) {
	ctx.Error(err).SetMeta(&ErrorMeta{Code: code})
	ctx.Abort()
}
