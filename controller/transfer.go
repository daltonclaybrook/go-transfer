package controller

import (
	"errors"
	"fmt"
	"github.com/daltonclaybrook/go-transfer/middle"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type Session struct {
	channel       chan []byte
	contentLength int
	timedOut      bool
}

type Transfer struct {
	sessions map[string]*Session
}

func NewTransfer() *Transfer {
	return &Transfer{make(map[string]*Session)}
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
func (transfer *Transfer) post(w http.ResponseWriter, r *http.Request, c middle.Context) {
	file := mux.Vars(r)["file"]
	ext := mux.Vars(r)["ext"]
	filename := fmt.Sprintf("%v.%v", file, ext)
	contentLength, err := strconv.Atoi(r.Header.Get("Content-Length"))
	fmt.Printf("posting file: %v\n", filename)

	if _, ok := transfer.sessions[filename]; ok {
		fmt.Fprintln(w, "this file is already being transfered")
	} else if (err != nil) || (contentLength <= 0) {
		fmt.Fprintln(w, "you must specify a content-length")
	} else {
		session := &Session{make(chan []byte, 2048), contentLength, false}
		transfer.sessions[filename] = session

		go closeChannelAfterDelay(session)
		err := performTransferRead(r, session, contentLength)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, err.Error())
		} else {
			fmt.Fprintln(w, "done")
		}
	}
}

// Used to receive a transfer
func (transfer *Transfer) get(w http.ResponseWriter, r *http.Request, c middle.Context) {
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

func performTransferRead(r *http.Request, session *Session, length int) error {
	totalRead := 0
	for {
		bytes := make([]byte, 1024)
		// fmt.Println("reading")
		count, err := r.Body.Read(bytes)

		if count > 0 {
			totalRead += count
			session.channel <- bytes[0:count]
		}

		if (err != nil) || (session.timedOut) {
			close(session.channel)
			break
		}
	}

	if totalRead < length {
		return errors.New("transfer could not be completed")
	}
	return nil
}

func (transfer *Transfer) performTransferWrite(w http.ResponseWriter, filename string, session *Session) {
	w.Header().Set("Content-Length", strconv.Itoa(session.contentLength))
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

func closeChannelAfterDelay(session *Session) {
	time.Sleep(time.Second * 120)
	session.timedOut = true
	<-session.channel
}
