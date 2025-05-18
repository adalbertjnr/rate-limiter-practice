package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	httpServerPort     string
	userLimiterStorage *UserLimitStorage
}

func (s *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("username")
	if user == "" {
		http.Error(w, "Missing query", http.StatusBadRequest)
		return
	}

	limiter := s.userLimiterStorage.CheckUser(user)

	if !limiter.Allow() {
		log.Printf("Request from user %s has been blocked", user)
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	log.Printf("Request from user %s has been allowed", user)
	fmt.Fprintf(w, "[%s] your request has been allowed!", user)
}

func main() {
	uls := &UserLimitStorage{
		m:      make(map[string]*RateLimiter),
		mu:     sync.RWMutex{},
		limit:  5,
		window: time.Second * 30,
	}

	srv := &Server{
		httpServerPort:     ":5000",
		userLimiterStorage: uls,
	}

	http.HandleFunc("GET /", srv.handleUser)

	log.Printf("Server starting on port %s", srv.httpServerPort)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
