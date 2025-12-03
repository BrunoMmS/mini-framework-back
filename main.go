package main

import (
	"strconv"
	"fmt"
	"encoding/json"
	"glac/router"
	"glac/server"
	"glac/custom_errors"

)
type User struct {
	ID int
	Name string
}
func main() {
	r := router.InitRouter()
	var users = []User{
	    {ID: 1, Name: "Juan"},
	    {ID: 2, Name: "Maria"},
	}

	r.Get("/", func(c *router.Context) {
		c.JSON(200, map[string]string{
			"message": "home",
		})
	})
	
	r.Get("/users/:id", func(c *router.Context) {
	    idStr := c.Params["id"]

	    id, err := strconv.Atoi(idStr)
	    if err != nil {
	        c.Error(custom_errors.BadRequest("id must be a number"))
	        return
	    }

	    if id < 1 || id > len(users) {
	        c.Error(custom_errors.NotFound("user"))
	        return
	    }

	    c.JSON(200, users[id-1])
	})

	r.Post("/", func(c *router.Context) {
    var payload struct {
        Message string `json:"message"`
    }

    if err := json.Unmarshal(c.Body, &payload); err != nil {
        fmt.Println("Invalid JSON:", err)
        return
    }

    fmt.Println("Message:", payload.Message)
})

	s := server.NewServer(r)
	s.Listen(":8080")
}
