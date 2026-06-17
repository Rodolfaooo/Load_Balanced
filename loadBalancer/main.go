type Backend struct {
	URL     string
	Alive   bool
	Latency time.Duration
}

var mu sync.RWMutex

func healthCheck(backends []*Backend) {

	for {
		for _, b := range backends {

			go func(b *Backend) {

				start := time.Now()

				resp, err := http.Get(b.URL + "/health")

				mu.Lock()
				defer mu.Unlock()

				if err != nil {
					b.Alive = false
					return
				}

				resp.Body.Close()

				b.Alive = true
				b.Latency = time.Since(start)

			}(b)
		}

		time.Sleep(3 * time.Second)
	}
}

func getLeastLatency(backends []*Backend) *Backend {

	mu.RLock()
	defer mu.RUnlock()

	var best *Backend

	for _, b := range backends {

		if !b.Alive {
			continue
		}

		if best == nil || b.Latency < best.Latency {
			best = b
		}
	}

	return best
}

func proxy(target string, w http.ResponseWriter, r *http.Request) {

	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ServeHTTP(w, r)
}

func handler(backends []*Backend) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		b := getLeastLatency(backends)

		if b == nil {
			http.Error(w, "Sem backends disponíveis", 503)
			return
		}

		fmt.Println("->", b.URL, "| latency:", b.Latency)

		proxy(b.URL, w, r)
	}
}

func main() {

	backends := []*Backend{
		{URL: "http://server1:8080"},
		{URL: "http://server2:8080"},
		{URL: "http://server3:8080"},
	}

	go healthCheck(backends)

	http.HandleFunc("/", handler(backends))

	fmt.Println("Load Balancer rodando em :8080")

	http.ListenAndServe(":8080", nil)
}