package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func getEnv(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func main() {
	port := getEnv("PORT", "8080")
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln("failed to connect to redis:", err)

	}
	fmt.Println("connected to redis:", pong)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		err := rdb.Incr(context.Background(), "hits").Err()
		if err != nil {
			log.Println(err)
			http.Error(w, "redis broke", http.StatusInternalServerError)
			return
		}
		val, err := rdb.Get(context.Background(), "hits").Int()
		if err != nil {
			log.Println(err)
			http.Error(w, "redis broke", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("this route was hit %d times", val)))
	})
	http.HandleFunc("GET /time", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().Format(time.RFC1123Z)))
	})
	fmt.Println("starting...")
	http.ListenAndServe(":"+port, nil)
}
