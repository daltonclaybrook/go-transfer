package controller

import (
	"fmt"
	"net/http"
)

type Transfer struct {
	channel       chan []byte
	contentType   string
	contentLength string
}

func (transfer *Transfer) Routes() []Route {
	crud := []CRUDRoute{
		CRUDRoute{Create, transfer.post},
		CRUDRoute{Find, transfer.get},
	}
	return RoutesFromCRUD("transfer", crud)
}

// Used to open a transfer session
func (transfer *Transfer) post(w http.ResponseWriter, r *http.Request) {
	transfer.channel = make(chan []byte)
	transfer.contentType = r.Header.Get("Content-Type")
	transfer.contentLength = r.Header.Get("Content-Length")

	// fmt.Println("received post..")

	for {
		bytes := make([]byte, 1024)
		count, err := r.Body.Read(bytes)
		// fmt.Println("read bytes..")

		if (err == nil) && (count > 0) {
			// fmt.Println("sending to channel")
			transfer.channel <- bytes[0:count]
			// fmt.Println("sent")
		} else {
			// fmt.Printf("error: %v\ncount: %v\n\n", err, count)
			close(transfer.channel)
			break
		}
	}

	fmt.Fprintln(w, "done")
}

// Used to receive a transfer
func (transfer *Transfer) get(w http.ResponseWriter, r *http.Request) {
	if transfer.contentLength != "" {
		w.Header().Set("Content-Length", transfer.contentLength)
	}

	w.Header().Set("Content-Type", "application/force-download")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", "minmus.mp4"))

	if transfer.channel != nil {
		w.WriteHeader(200)
		for {
			bytes, ok := <-transfer.channel
			if ok {
				fmt.Println("writing...")
				w.Write(bytes)
				// w.(http.Flusher).Flush()
			} else {
				fmt.Println("done")
				transfer.channel = nil
				break
			}
		}
	} else {
		fmt.Fprintln(w, "did not find an open channel")
	}
}
