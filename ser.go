package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var addr = flag.String("addr", "localhost", "address for listening")
var port = flag.String("p", "8080", "port for listening")

func usage() {
	fmt.Fprintf(os.Stderr, `ser is a simple http server.

Command line usage:

$ ser [--help] [-addr <addr>] [-p <port>] [<dir>]

options:
  --help
        print this message
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
examples:
ser
ser .
	serve . directory using port 8080
ser -p 8088
	serve . directory using port 8088
ser ..
	serve .. directory using port 8080
ser -p 8088 ..
	serve .. directory using port 8088
ser -addr 0.0.0.0 -p 9999
		serve . directory using address 0.0.0.0:9999
`)
	os.Exit(2)
}

func main() {
	log.SetPrefix("ser: ")
	log.SetFlags(log.Lmsgprefix | log.LstdFlags | log.Lshortfile)
	flag.Usage = usage
	flag.Parse()
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

	n, err := strconv.Atoi(*port)
	if err != nil || n <= 1000 {
		l.Printf("invalid port number %d: %v", n, err)
		return
	}

	dir := "."

	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
	}

	logger := logging(l)

	path, err := filepath.Abs(dir)
	if err != nil {
		l.Printf("failed to get absolute path of %s: %v", dir, err)
		flag.Usage()
		return
	}

	_, err = os.Lstat(path)
	if os.IsNotExist(err) {
		l.Printf("folder %s does not exist!", path)
		return
	}

	listen := *addr + ":" + *port
	l.Printf("serving %s at %s", path, listen)
	http.Handle("/", logger(noCache(http.FileServer(http.Dir(path)))))
	log.Fatal(http.ListenAndServe(listen, nil))
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r, r.Method, r.URL.Path)
			}()
			next.ServeHTTP(w, r)
		})
	}
}

var noCacheHeaders = map[string]string{
	"Expires":         time.Unix(0, 0).Format(time.RFC1123),
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func noCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
