package pkg

import (
	"github.com/winstonjr/goexpert-desafio-rate-limiter/internal/entity"
	"log"
	"net/http"
	"strings"
)

func TokenBucket(store entity.FilterStoreInterface) func(http.Handler) http.Handler {
	filter := NewFilter(store)
	return filter.TokenBucketHandler
}

func (f *Filter) TokenBucketHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("TokenBucket Called: %s", r.URL.Path)
		key := r.Header.Get("API_KEY")
		if strings.TrimSpace(key) == "" {
			key = getIPAddress(r)
		}

		if f.store.LimitExceeded(key) {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame\n", 429)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getIPAddress(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
