package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Session struct {
	channel       chan []byte
	contentLength string
}

type Transfer struct {
	sessions map[string]Session
}

func NewTransfer() *Transfer {
	return &Transfer{make(map[string]Session)}
}

func (transfer *Transfer) Routes() []Route {
	return []Route{
		Route{"/{file}.{ext}", []Handler{
			Handler{"post", transfer.post},
			Handler{"get", transfer.get},
		}},
	}
}

// Used to open a transfer session
func (transfer *Transfer) post(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	ext := mux.Vars(r)["ext"]
	filename := fmt.Sprintf("%v.%v", file, ext)
	contentLength := r.Header.Get("Content-Length")
	fmt.Printf("posting file: %v\n", filename)

	if _, ok := transfer.sessions[filename]; ok {
		fmt.Fprintln(w, "this file is already being transfered")
	} else if contentLength == "" {
		fmt.Fprintln(w, "you must specify a content-length")
	} else {
		session := Session{make(chan []byte), contentLength}
		transfer.sessions[filename] = session
		performTransferRead(r, session)
		fmt.Fprintln(w, "done")
	}
}

// Used to receive a transfer
func (transfer *Transfer) get(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	ext := mux.Vars(r)["ext"]
	filename := fmt.Sprintf("%v.%v", file, ext)
	fmt.Printf("getting file: %v\n", filename)

	if session, ok := transfer.sessions[filename]; ok {
		transfer.performTransferWrite(w, filename, session)
	} else {
		fmt.Fprintln(w, "A transfer has not begun at this endpoint.")
	}
}

/*
	Helper methods
*/

func performTransferRead(r *http.Request, session Session) {
	for {
		bytes := make([]byte, 1024)
		count, err := r.Body.Read(bytes)
		// fmt.Println("read bytes..")

		if (err == nil) && (count > 0) {
			// fmt.Println("sending to channel")
			session.channel <- bytes[0:count]
			// fmt.Println("sent")
		} else {
			// fmt.Printf("error: %v\ncount: %v\n\n", err, count)
			close(session.channel)
			break
		}
	}
}

func (transfer *Transfer) performTransferWrite(w http.ResponseWriter, filename string, session Session) {
	w.Header().Set("Content-Length", session.contentLength)
	w.Header().Set("Content-Type", "application/force-download")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", filename))

	for {
		bytes, ok := <-session.channel
		if ok {
			w.Write(bytes)
		} else {
			delete(transfer.sessions, filename)
			break
		}
	}
}
