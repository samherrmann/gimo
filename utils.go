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
