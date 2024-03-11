package webapi

import (
	"fmt"
	"log"
	http "net/http"
)

func StartWebAPI(port int) {
	http.HandleFunc("/api/image", handleImage)

	address := "127.0.0.1:" + fmt.Sprint(port)

	log.Println("web server starting @", "http://"+address)
	log.Fatal(http.ListenAndServe(address, nil))
}
