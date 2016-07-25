# Gimo
A Go library to build CRUD APIs with [Gin](https://github.com/gin-gonic/gin) and [mgo](https://github.com/go-mgo/mgo/tree/v2) (mongoDB).

<span style="color: red">This project is a work in progress!</span>

[API Documentation](https://godoc.org/github.com/samherrmann/gimo)

## Features
* Add basic CRUD resources to your app quickly without the need to write any handler functions.
* Add business logic with middleware.
* Gin and mgo instances used by Gimo are accessible outside of Gimo, allowing you to handle edge cases that fall outside of Gimo's capabilities.
* Gimo does not use reflection.

## Steps to add a resource
1. Define the resource's data model.
2. Give your resource a name. The name is used to generate the resource's URL(s) and its database collection name.
3. Declare which actions are supported on your resource, i.e. Create, read, update, delete, and list.
4. (If required) Add middleware to handle specific business logic.

## Installation
```sh
go get github.com/samherrmann/gimo
```

## Examples

The examples assumes that you are familiar with Gin and mgo.

### Simple resource

```go
func main() {
    // Create a Gin router.
    router := gin.Default()
    group := router.Group("/v1")

    // Define the mgo dial info.
    info := &mgo.DialInfo{}
    info.Addrs = []string{"127.0.0.1"}
    info.Database = "my-store"
    info.Timeout = 2 * time.Second

    // Create a new Gimo library by passing it
    // the Gin router and the mgo dial info.
    lib := gimo.Default(group, info)
    defer lib.Terminate()

    // Add a "books" resource with all actions enabled.
    // The book data model is defined below.
    res := lib.Resource("books", &models.Book{})
    res.Create()
    res.Read()
    res.Update()
    res.Delete()
    res.List()

    // Start the server.
    router.Run(":8080")
}

// Define the book data model.
// The data model struct must implement Gimo's
// Document interface, i.e. it must contain a 
// "SetID", "GetID", "New", and "Slice" method.
// The "SetID" and "GetID" methods are provided 
// through Gimo's embedded "DocumentBase" struct.
type Book struct {
    gimo.DocumentBase `json:",inline"    bson:",inline"`
    Title             string `json:"title"`
    Author            string `json:"author"`
    Publisher         string `json:"publisher"`
}

func (b *Book) New() gimo.Document {
    return &Book{}
}

func (b *Book) Slice() interface{} {
    return &[]Book{}
}
```

The endpoints that Gimo created in this example are as follows:

    Create:     POST    /v1/books
    Read:       GET     /v1/books/:id
    Update:     PUT     /v1/books/:id
    Delete:     DELETE  /v1/books/:id
    List:       GET     /v1/books

This example demonstrated how to add a CRUD resource without having to write a single Gin handler function to interact with mongoDB.

### Using middleware

Add Gin middleware to any of the action methods the same way you would in Gin directly.

```go
...
res.Create(myGinMiddleware /*, anotherMiddleware, ... */ )
res.Read( /* middlewareX, middlewareY, ... */ )
res.Update( /* middlewareZ, middlewareA, ... */ )
res.Delete( /* middlewareB, middlewareC, ... */ )
res.List( /* middlewareFoo, middlewareBar, ... */)
...

func myGinMiddleware(ctx *gin.Context) {
    ...
}
```
Gimo merges your middleware between a couple of its own middleware that it uses underneath the hood. The chain that Gimo generates looks like this...

```
handleErrors, serializeResponse, parseRequest, yourMiddleware, ..., mongoDBHandler
```

Note that the first two middleware (`handleErrors` and `serializeResponse`) only "do stuff" after the request, not before. If you are unsure about what that means, see the Gin example for [Custom Middleware](https://github.com/gin-gonic/gin#custom-middleware).

#### Gin context key/value pairs

When Gimo's `parseRequest` middleware parses the incoming data, for the `Create` and `Update` method, it stores the parsed data in a `Document` under a context key named `request`. You are therefore able to access the parsed data in your middleware as follows:

```go
func myGinMiddleware(ctx *gin.Context) {
    doc, exists := ctx.Get("request")
    ...
}
```
Similarly, the `mongoDBHandler` sets the response from the mongoDB database in a Gin context key named `response`. A middleware that accesses the response may look as follows:

```go
func myGinMiddleware(ctx *gin.Context) {
    ctx.Next()
    doc, exists := ctx.Get("response")
    ...
}
```
Note that the response is accessed after calling `ctx.Next()`! Accessing the response key before calling `ctx.Next()` would result in `exists` equating to `false`.

##### Custom request/response context keys
The Gin context key names `request` and `response` are what Gimo creates/uses when the Gimo library instance is created with the `Default` function:

```go
lib := gimo.Default(routerGroup, mgoInfo)
```

If you need to change the key names from their default names, create the Gin library instance with the `New` function instead. The `New` function allows you to set the custom context keys:

```go
lib := gimo.New(routerGroup, mgoInfo, "req", "res")
```

Of course if you assign custom keys, the keys used by your middleware need to match. Therefore the following example demonstrates what may be a better way to write your middleware to ensure that they always use the correct key names:

```go
func myGinMiddleware(r *gimo.Resource) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        ctx.Next()
        doc, exists := ctx.Get(r.ResponseCtxKey)
        ...
    }
}
```