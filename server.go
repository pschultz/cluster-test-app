package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	flag.Parse()
	port := flag.Arg(0)
	if port == "" {
		port = "80"
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out := io.MultiWriter(w, os.Stdout)

		fmt.Fprintf(out, "%s %s %s\n", r.Method, r.URL.String(), r.Proto)
		fmt.Fprintf(out, "Host: %s\n", r.Host)
		r.Header.Write(out)

		fmt.Fprintf(w, "\n%s:%s\n", hostname, port)
	})

	fmt.Fprintf(os.Stdout, "Listening on 0.0.0.0:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
