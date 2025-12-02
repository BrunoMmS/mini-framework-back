package main

import (
    "fmt"
    "glac/router"
    "glac/server"
)

func main() {
    r := router.InitRouter()

    r.Handle("GET", "/hello", func() (any, error) {
        return "Hola desde GET /hello", nil
    })

    r.Handle("POST", "/register", func() (any, error) {
        return map[string]string{"status": "ok"}, nil
    })

    s := server.NewServer("tcp", ":8080", r)

    fmt.Println("Servidor corriendo en http://localhost:8080")
    s.Listen()
}
