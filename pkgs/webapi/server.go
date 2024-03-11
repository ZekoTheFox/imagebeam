package webapi

import (
	"fmt"
	"log"
	http "net/http"
)

type Image struct {
	Url string
}

var Images chan Image = make(chan Image)
var imageCache chan []byte = make(chan []byte)

func StartWebAPI(port int) {
	http.HandleFunc("GET /image", handleImage)

	address := "127.0.0.1:" + fmt.Sprint(port)

	log.Println("web server starting @", "http://"+address)
	log.Fatal(http.ListenAndServe(address, nil))
}
