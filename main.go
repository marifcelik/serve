package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"
)

const version = "0.1"

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
	var showVersion bool

	flag.StringVar(&host, "host", "localhost", "host to listen on")
	flag.StringVar(&host, "h", "localhost", "host to listen on")
	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.StringVar(&port, "p", "8080", "port to listen on")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")

	flag.Usage = func() {
		usage := `Usage: serve [options] <path>

Options:
  -h, --host string
      host to listen on (default "localhost")
  -p, --port string
      port to listen on (default "8080")

Arguments:
  <path>  path to serve (default current working directory)`

		fmt.Println(usage)
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("v%s\n", version)
		os.Exit(0)
	}

	var err error
	path, err = os.Getwd()
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	} else if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", logger(fs))

Serve:
	addr := net.JoinHostPort(host, port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		if operr, ok := err.(*net.OpError); ok && errors.Is(operr.Err, syscall.EADDRINUSE) {
			goto ChangePort
		}
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Printf("Serving %s on %s\n", path, addr)
	log.Fatal(http.Serve(l, nil))

ChangePort:
	newPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}
	port = strconv.Itoa(newPort + 1)
	log.Printf("WARN Address already in use, changing port to %s\n", port)
	goto Serve
}
