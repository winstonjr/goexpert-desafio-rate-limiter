package limiter

func TokenBucket(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("passou por aqui ;D: %s", r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
