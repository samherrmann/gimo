package main

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
	"github.com/samherrmann/gimo"
	"github.com/samherrmann/gimo/example/models"
)

func main() {

	// Initialize Gin
	router := gin.Default()
	group := router.Group("/v1")

	// Define the mgo configuration
	info := &mgo.DialInfo{}
	info.Addrs = []string{"localhost:27017"}
	info.Database = "my-store"
	info.Timeout = 2 * time.Second

	// Create a new gimo library
	lib := gimo.Default(group, info)
	defer lib.Terminate()

	// Gimo's "Default" function above dialed the mongoDB server/cluster and established
	// a session. You have access to the establisehd session and can make changes
	// as you wish...
	lib.Session.SetMode(mgo.Monotonic, true)

	// Add a "books" resource with all actions enabled
	res := lib.Resource("books", &models.Book{})
	res.Create( /* you can add middleware in here to handle your business logic */ )
	res.Read( /* middleware */ )
	res.Update( /* middleware */ )
	res.Delete( /* middleware */ )
	res.List( /* middleware */ )

	router.Run(":8080")
}
