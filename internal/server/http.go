package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	// take in address, returns server
	// wrap with http.Server so user just calls ListenAndServe
	httpsrv := newHTTPServer()
	r := mux.NewRouter()
	// gorilla handles routing requests
	r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

type ProduceRequest struct {
	// caller wants to append to log
	Record Record `json:"record"`
}

type ProduceResponse struct {
	// server tells caller what offset the record was stored at
	Offset uint64 `json:"offset"`
}
type ConsumeRequest struct {
	// caller wants to read some logs at offset
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	// server returns record
	Record Record `json:"record"`
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	// handle Produce:
	// 1. try decoding request.json
	// 2. try main logic (append record)
	// 3. try sending response back to caller
	// if any of these fail, throw an error (error check after each one)
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	// handle Consume:
	// 1. try decoding request.json
	// 2. try main logic (read record)
	// 3. try sending response back to caller
	// if any of these fail, throw an error (error check after each one)
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	record, err := s.Log.Read(req.Offset) //offset this time
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
