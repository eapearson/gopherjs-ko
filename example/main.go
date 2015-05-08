package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

var assetRegexp = regexp.MustCompile(`\.(html|jpg|png|ico|css|js|json|eot|svg|ttf|woff|woff2)$`)

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

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		val := strconv.Itoa(rand.Int())
		w.Write([]byte(`{"data": "Some data from the backend", "value": ` + val + `}`))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if assetRegexp.MatchString(r.URL.Path) {
			http.ServeFile(w, r, "static"+r.URL.Path)
		} else {
			http.ServeFile(w, r, "static/index.html")
		}
	})

	log.Println("Listening on :8000 ...")
	http.ListenAndServe(":8000", nil)
}
