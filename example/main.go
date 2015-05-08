package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	go func() {
		log.Println("Running gopherjs ...")
		cmd := exec.Command("gopherjs", "build",
			"github.com/mibitzi/gopherjs-ko/example/js",
			"-o", "static/demo.js")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}()

	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Listening on :8000 ...")
	http.ListenAndServe(":8000", nil)
}
