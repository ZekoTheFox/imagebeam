package webapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handleImage(writer http.ResponseWriter, request *http.Request) {
	select {
	case image := <-Images:
		download, err := http.Get(image.Url)
		if err != nil || download.StatusCode != 200 {
			log.Println("warning: failed to download media (status", fmt.Sprint(download.StatusCode)+")")
			io.WriteString(writer, "")
			return
		}

		imageData, err := io.ReadAll(download.Body)
		if err != nil {
			log.Println("warning: failed to read link data")
			io.WriteString(writer, "")
			return
		}

		writer.Write(imageData)
	default:
		io.WriteString(writer, "") // respond with nothing
	}
}
