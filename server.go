package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

var hostname, port string

func main() {
	flag.Parse()
	port = flag.Arg(0)
	if port == "" {
		port = "80"
	}
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	codes := []int{
		200, 201, 202,
		400, 401, 403, 404, 410, 418,
		500, 501, 502, 503, 504,
	}
	for _, c := range codes {
		http.Handle(fmt.Sprintf("/%d/", c), statusHandler(c))
		http.Handle(fmt.Sprintf("/%d", c), statusHandler(c))
	}
	http.Handle("/sleep/", sleepHandler())
	http.Handle("/sleep", sleepHandler())
	http.Handle("/", statusHandler(200))

	fmt.Fprintf(os.Stdout, "Listening on 0.0.0.0:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func statusHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, false)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(status)
		out := io.MultiWriter(w, os.Stdout)
		out.Write(b)

		fmt.Fprintf(w, "\n%s:%s\n", hostname, port)
	}
}

func sleepHandler() http.HandlerFunc {
	h := statusHandler(200)

	return func(w http.ResponseWriter, r *http.Request) {
		durationStr := strings.Trim(r.URL.Path[len("/sleep"):], "/")
		if durationStr == "" {
			time.Sleep(time.Second)
			h(w, r)
			return
		}

		d, err := time.ParseDuration(durationStr)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		time.Sleep(d)
		h(w, r)
	}
}
