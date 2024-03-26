package webapi

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	http "net/http"
)

type Image struct {
	Url string
}

var (
	//go:embed overlay/index.html
	overlayFile string

	Images chan Image = make(chan Image)
)

func StartWebAPI(port int) {
	http.HandleFunc("GET /image", handleImage)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		io.WriteString(w, overlayFile)
	})

	address := "127.0.0.1:" + fmt.Sprint(port)

	log.Println("web server starting @", "http://"+address)
	log.Fatal(http.ListenAndServe(address, nil))
}
