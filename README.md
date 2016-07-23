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

## Example

This example assumes that you are familiar with Gin and mgo.

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