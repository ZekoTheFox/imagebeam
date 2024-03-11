package webapi

import (
	"io"
	"net/http"
)

func handleImage(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "https://picsum.photos/320/180") // placeholder image
}
