package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/winstonjr/goexpert-desafio-rate-limiter/configs"
	"github.com/winstonjr/goexpert-desafio-rate-limiter/internal/entity"
	"github.com/winstonjr/goexpert-desafio-rate-limiter/internal/infra/database"
	"github.com/winstonjr/goexpert-desafio-rate-limiter/pkg"
	"log"
	"net/http"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	fsi := createStoreStrategy(config)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "/*"))
	r.Route("/", func(r chi.Router) {
		r.Use(pkg.TokenBucket(fsi))
		r.Get("/", helloWorldHandler)
	})
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("helloWorldHandler called")
	_, err := w.Write([]byte("helloWorldHandler"))
	if err != nil {
		log.Println(err)
	}
}

func createStoreStrategy(config *configs.Conf) entity.FilterStoreInterface {
	if config.StoreKind == "redis" {
		ctx := context.Background()
		rdb := redis.NewClient(&redis.Options{
			Addr:     config.RedisAddress,
			Password: config.RedisPassword,
			DB:       config.RedisDb,
		})
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			log.Fatal("could not ping redis", err)
		}

		fsi, err := database.NewFilterStoreRedis(config.RateLimiterRules, rdb)
		if err != nil {
			log.Fatal(err)
		}

		return fsi
	}

	fsi, err := database.NewFilterStoreInMemory(config.RateLimiterRules)
	if err != nil {
		log.Fatal(err)
	}
	return fsi
}
