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

func StartWebAPI(port int) {
	http.HandleFunc("GET /images", handleImage)

	address := "127.0.0.1:" + fmt.Sprint(port)

	log.Println("web server starting @", "http://"+address)
	log.Fatal(http.ListenAndServe(address, nil))
}
