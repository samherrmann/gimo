package gimo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// DefaultRequestCtxKey is the default Gin context key uder
	// which the parsed request body is stored. The request context
	// key can be customized with the "New"" function.
	DefaultRequestCtxKey = "request"

	// DefaultResponseCtxKey is the default Gin context key uder
	// which the response struct is stored. The response context key
	// can be customized with the "New"" function.
	DefaultResponseCtxKey = "response"

	// The path parameter key for a resource ID.
	idPathParamKey = "id"
)

type (
	// Library is Gimo's core instance. It contains settings
	// that are common to all resources.
	Library struct {
		BaseGroup      *gin.RouterGroup
		Session        *mgo.Session
		RequestCtxKey  string
		ResponseCtxKey string
	}

	// Resource represents a single CRUD resource.
	Resource struct {
		*Library
		Name  string
		Group *gin.RouterGroup
		Doc   Document
	}
)

// Default returns a Library instance with default internal settings.
func Default(baseGroup *gin.RouterGroup, dbInfo *mgo.DialInfo) *Library {
	return New(baseGroup, dbInfo, DefaultRequestCtxKey, DefaultResponseCtxKey)
}

// New returns a Library instance.
func New(baseGroup *gin.RouterGroup, dbInfo *mgo.DialInfo, requestCtxKey string, responseCtxKey string) *Library {
	if requestCtxKey == "" {
		requestCtxKey = DefaultRequestCtxKey
	}
	if responseCtxKey == "" {
		responseCtxKey = DefaultResponseCtxKey
	}

	return &Library{
		BaseGroup:      baseGroup,
		Session:        dialDB(dbInfo),
		RequestCtxKey:  requestCtxKey,
		ResponseCtxKey: responseCtxKey,
	}
}

// Resource returns a Resource instance.
func (lib *Library) Resource(name string, doc Document) *Resource {
	group := lib.BaseGroup.Group(name)

	r := &Resource{
		Library: lib,
		Name:    name,
		Group:   group,
		Doc:     doc,
	}

	group.Use(r.handleErrors)
	group.Use(r.serializeResponse)
	return r
}

// Terminate closes the mongoDB session.
func (lib *Library) Terminate() {
	lib.Session.Close()
}

// Create adds a Gin handler function that allows
// a client to create a new document in the mongoDB
// collection.
func (r *Resource) Create(mw ...gin.HandlerFunc) {
	h := func(ctx *gin.Context) {
		c := r.Session.Clone().DB("").C(r.Name)
		defer c.Database.Session.Close()

		doc := ctx.MustGet(r.RequestCtxKey).(Document)
		doc.SetID(bson.NewObjectId().Hex())
		err := c.Insert(doc)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Set(r.ResponseCtxKey, doc)
	}

	chain := append([]gin.HandlerFunc{r.parseRequest}, mw...)
	chain = append(chain, h)
	r.Group.POST("/", chain...)
}

// Read adds a Gin handler function that allows
// a client to get a single document from the
// mongoDB collection.
func (r *Resource) Read(mw ...gin.HandlerFunc) {
	h := func(ctx *gin.Context) {
		c := r.Session.Clone().DB("").C(r.Name)
		defer c.Database.Session.Close()

		doc := r.Doc.New()
		err := c.FindId(ctx.Param(idPathParamKey)).One(doc)
		if err == mgo.ErrNotFound {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Set(r.ResponseCtxKey, doc)
	}

	chain := append(mw, h)
	r.Group.GET("/:"+idPathParamKey, chain...)
}

// Update adds a Gin handler function that allows
// a client to update an existing document in the
// mongoDB collection.
func (r *Resource) Update(mw ...gin.HandlerFunc) {
	h := func(ctx *gin.Context) {
		c := r.Session.Clone().DB("").C(r.Name)
		defer c.Database.Session.Close()

		doc := ctx.MustGet(r.RequestCtxKey).(Document)
		doc.SetID(ctx.Param(idPathParamKey))
		err := c.UpdateId(doc.GetID(), bson.M{"$set": doc})
		if err == mgo.ErrNotFound {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Set(r.ResponseCtxKey, doc)
	}

	chain := append([]gin.HandlerFunc{r.parseRequest}, mw...)
	chain = append(chain, h)
	r.Group.PUT("/:"+idPathParamKey, chain...)
}

// Delete adds a Gin handler function that allows
// a client to remove a document from the mongoDB
// collection.
func (r *Resource) Delete(mw ...gin.HandlerFunc) {
	h := func(ctx *gin.Context) {
		c := r.Session.Clone().DB("").C(r.Name)
		defer c.Database.Session.Close()

		err := c.RemoveId(ctx.Param(idPathParamKey))
		if err == mgo.ErrNotFound {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Set(r.ResponseCtxKey, nil)
	}

	chain := append(mw, h)
	r.Group.DELETE("/:"+idPathParamKey, chain...)
}

// List adds a Gin handler function that allows
// a client to get all documents from the mongoDB
// collection.
func (r *Resource) List(mw ...gin.HandlerFunc) {
	h := func(ctx *gin.Context) {
		c := r.Session.Clone().DB("").C(r.Name)
		defer c.Database.Session.Close()

		docs := r.Doc.Slice()
		err := c.Find(nil).All(docs)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Set(r.ResponseCtxKey, docs)
	}

	chain := append(mw, h)
	r.Group.GET("/", chain...)
}

// parseRequest is a Gin handler function that parses the
// JSON in the request body and stores the parsed result
// in the Gin context.
func (r *Resource) parseRequest(ctx *gin.Context) {
	doc := r.Doc.New()
	err := ctx.BindJSON(doc)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.Set(r.RequestCtxKey, doc)
	ctx.Next()
}

// serializeResponse is a Gin handler function that serializes
// the struct stored in the "response" Gin context to JSON and
// writes it to the response body.
func (r *Resource) serializeResponse(ctx *gin.Context) {
	ctx.Next()

	if len(ctx.Errors) > 0 {
		return
	}

	doc, exists := ctx.Get(r.ResponseCtxKey)
	if !exists || doc == nil {
		ctx.Status(http.StatusNoContent)
		return
	}
	ctx.JSON(http.StatusOK, doc)
}

// handleErrors checks if any errors have occured in the request
// chain. If an error has occured but no status code was written
// to the response, it writes a status code of 500 (Internal server
// error). Note that all client errors should be handled in
// middleware by calling Gin's 'AbortWithError' function. Also, if
// a middleware wishes to expose to the client the error message
// that it passes to the 'AbortWithError' function, then Gin's
// 'SetType(gin.ErrorTypePublic)' function should be called. It is
// assumed that if 'AbortWithError' is called, that this is always
// the first error in the list of possible error that may occur
// within a request chain.
func (r *Resource) handleErrors(ctx *gin.Context) {
	ctx.Next()

	if len(ctx.Errors) == 0 {
		return
	}

	if !ctx.Writer.Written() {
		ctx.Status(http.StatusInternalServerError)
	}

	var errMsg string
	if ctx.Errors[0].Type == gin.ErrorTypePublic {
		errMsg = ctx.Errors[0].Error()

	} else {
		errMsg = http.StatusText(ctx.Writer.Status())
	}
	ctx.Writer.WriteString(errMsg)
}
