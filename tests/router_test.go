package main

import (
	"encoding/json"
	"errors"
	"glac/router"
	"glac/server"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func readJSON[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	var result T
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
	return result
}

// TEST 1 — GET /
func TestRouteHome(t *testing.T) {
	r := router.InitRouter()

	r.Get("/", func(c *router.Context) {
		c.JSON(200, map[string]string{"message": "home"})
	})

	srv := server.NewServer(r)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	resp := readJSON[map[string]string](t, w)

	if resp["message"] != "home" {
		t.Fatalf("expected 'home', got %s", resp["message"])
	}
}

// TEST 2 — GET /users/:id
func TestRouteDynamic(t *testing.T) {
	users := []string{"Ana", "Luis", "Marta"}

	r := router.InitRouter()

	r.Get("/users/:id", func(c *router.Context) {
		id, _ := strconv.Atoi(c.Params["id"])
		if id < 1 || id > len(users) {
			c.JSON(404, map[string]string{"error": "not found"})
			return
		}
		c.JSON(200, map[string]string{"name": users[id-1]})
	})

	srv := server.NewServer(r)
	req := httptest.NewRequest(http.MethodGet, "/users/2", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	resp := readJSON[map[string]string](t, w)

	if resp["name"] != "Luis" {
		t.Fatalf("expected Luis, got %s", resp["name"])
	}
}

// TEST 3 — Middleware modifies context
func TestMiddlewareModify(t *testing.T) {

	logMW := func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			c.Params["debug"] = "1"
			next(c)
		}
	}

	r := router.InitRouter()

	r.Get("/test", func(c *router.Context) {
		c.JSON(200, map[string]string{"debug": c.Params["debug"]})
	}, logMW)

	srv := server.NewServer(r)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	resp := readJSON[map[string]string](t, w)

	if resp["debug"] != "1" {
		t.Fatalf("middleware did not modify params")
	}
}

// TEST 4 — Middleware blocks the request
func TestMiddlewareBlock(t *testing.T) {

	blocker := func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			c.JSON(401, map[string]string{"error": "blocked"})
			return
		}
	}

	r := router.InitRouter()

	r.Get("/private", func(c *router.Context) {
		c.JSON(200, map[string]string{"ok": "should-not-happen"})
	}, blocker)

	srv := server.NewServer(r)
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// TEST 5 — 404 Not Found
func TestNotFound(t *testing.T) {
	r := router.InitRouter()
	srv := server.NewServer(r)

	req := httptest.NewRequest(http.MethodGet, "/not-exists", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// TEST 6 — Panic → Internal Server Error 500
func TestInternalServerError(t *testing.T) {
	r := router.InitRouter()

	r.Get("/panic", func(c *router.Context) {
		panic(errors.New("boom"))
	})

	srv := server.NewServer(r)
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// Benchmark (performance)
func BenchmarkRouter(b *testing.B) {
	r := router.InitRouter()

	r.Get("/bench", func(c *router.Context) {
		c.JSON(200, map[string]string{"ok": "bench"})
	})

	srv := server.NewServer(r)

	req := httptest.NewRequest(http.MethodGet, "/bench", nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
	}
}
