package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

var host, port, path string

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("- [%s] %s \"%d\" %s", r.RemoteAddr, r.Method, rw.statusCode, r.URL.String())
	})
}

func init() {
	flag.StringVar(&host, "host", "localhost", "host to listen on")
	flag.StringVar(&host, "h", "localhost", "host to listen on")
	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.StringVar(&port, "p", "8080", "port to listen on")

	flag.Usage = func() {
		usage := `Usage: serve [options] <path>

Options:
  -h, --host string
      host to listen on (default "localhost")
  -p, --port string
      port to listen on (default "8080")

Arguments:
  <path>  path to serve (default ".")`

		fmt.Println(usage)
	}

	flag.Parse()

	path = "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}
}

func main() {
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", logger(fs))

	addr := net.JoinHostPort(host, port)
	fmt.Printf("Serving %s on %s\n", path, addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
