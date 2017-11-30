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
	port := ":20001"
	go func() {
		fs := http.FileServer(http.Dir("static"))
		http.Handle("/", fs)
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
