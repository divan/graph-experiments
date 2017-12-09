package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

func startWeb() error {
	port := ":20002"
	go func() {
		fs := http.FileServer(http.Dir("web"))
		http.Handle("/", noCacheMiddleware(fs))
		log.Fatal(http.ListenAndServe(port, nil))
	}()
	time.Sleep(1 * time.Second)
	startBrowser("http://localhost" + port)
	select {}
}

// startBrowser tries to open the URL in a browser
// and reports whether it succeeds.
//
// Orig. code: golang.org/x/tools/cmd/cover/html.go
func startBrowser(url string) error {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	fmt.Println("If browser window didn't appear, please go to this url:", url)
	return cmd.Start()
}

func noCacheMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=0, no-cache")
		h.ServeHTTP(w, r)
	})
}
