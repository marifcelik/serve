package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var fs = flag.NewFlagSet("serve", flag.ContinueOnError)
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
	fs.StringVar(&host, "host", "localhost", "host to listen on")
	fs.StringVar(&host, "h", "localhost", "host to listen on")
	fs.StringVar(&port, "port", "8080", "port to listen on")
	fs.StringVar(&port, "p", "8080", "port to listen on")
	fs.StringVar(&path, "path", ".", "path to serve")
	fs.StringVar(&path, "P", ".", "path to serve")

	fs.Usage = func() {
		usage := `Usage: serve [options]

Options:
  -h, --host string
	  host to listen on (default "localhost")
  -p, --port string
	  port to listen on (default "8080")
  -P, --path string
	  path to serve (default ".")`

		fmt.Println(usage)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("error parsing flags: %v", err)
	}
}

func main() {
	server := http.FileServer(http.Dir(path))
	http.Handle("/", logger(server))

	addr := net.JoinHostPort(host, port)
	fmt.Printf("Serving %s on %s\n", path, addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
