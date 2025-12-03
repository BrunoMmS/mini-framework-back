package main

import (
	"fmt"
	"glac/router"
	"glac/server"
	"strconv"
)

func getAllUsers(users map[int]string) (map[int]string) {
    return users
}
func main() {
    r := router.InitRouter()

    r.Handle("GET", "/hello", func(c *router.Context){
        c.Text(200, "Hola desde GET /hello")
    })

    r.Handle("POST", "/register", func(c *router.Context) {
        c.JSON(200, "Registro OK")
    })
    
    users := make(map[int]string)
    users[1] = "ASD"
    users[2] = "DSA"    
    r.Handle("GET", "/users/:id", func(c *router.Context) {
        id := c.Params["id"]
        idInt, _ := strconv.Atoi(id)
        user, ok := users[idInt]
        if !ok{
            c.Error(404, "User not found")
            return
        }
        c.JSON(200, map[string]string{"user": user})
    })
    

    r.Handle("GET", "/users", func(c *router.Context) {
        allUsers := getAllUsers(users)  
        c.JSON(200, allUsers)
    } )

    s := server.NewServer("tcp", ":8080", r)

    fmt.Println("Servidor corriendo en http://localhost:8080")
    s.Listen()
}
