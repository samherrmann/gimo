# Gimo
A Go library to build CRUD APIs with [Gin](https://github.com/gin-gonic/gin) and [mgo](https://github.com/go-mgo/mgo/tree/v2).

[Gimo GoDoc](https://godoc.org/github.com/samherrmann/gimo)

<span style="color: red">This project is a work in progress!</span>

## What problem is Gimo solving?
Not having to write separate endpoint handlers for every resource. For common use cases, it should be possible to break down the process of adding a resource into the following steps:

1. Define the resource's data model
2. Define the resource's name (to generate the resource's URL(s) and its database collection)
3. Define the supported actions, i.e. Create, read, update, delete, and list
4. (If required) Define the middleware to handle specific business logic

A developer should not have to write endpoint handlers for every resource that moves data in and out of the database.

## Gimo design goals
* To not hide Gin or mgo underneath Gimo's hood. This allows you, the developer, to not be locked into yet another library. You have the ability to use Gimo where you see fit and cover edge cases that fall outside of Gimo's capabilities by just using Gin and mgo.
* To not use reflection
* To allow you to handle business logic with middleware

## Installation
```sh
go get github.com/samherrmann/gimo
```

## Example

```go
func main() {
    // Create a Gin router.
    router := gin.Default()
    group := router.Group("/v1")

    // Define the mgo dial info.
    info := &mgo.DialInfo{}
    info.Addrs = []string{"localhost:27017"}
    info.Database = "my-store"
    info.Timeout = 2 * time.Second

    // Create a new Gimo library by passing it
    // the Gin router and the mgo dial info.
    lib := gimo.Default(group, info)
    defer lib.Terminate()

    // Add a "books" resource with all actions enabled.
    res := lib.Resource("books", &models.Book{})
    res.Create()
    res.Read()
    res.Update()
    res.Delete()
    res.List()

    // Start the server.
    router.Run(":8080")
}

// Define the data model for the book resource.
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

In this example, Gimo performed the following two functions:

1. It created the following paths:

    ```
    Create:     POST    /v1/books
    Read:       GET     /v1/books/:id
    Update:     PUT     /v1/books/:id
    Delete:     DELETE  /v1/books/:id
    List:       GET     /v1/books
    ```

2. It attached handler functions that handle the data flow between the HTTP requests and the database.

Note: Gimo used the name `books` to generate the paths, but it also used it as the mongoDB collection name.