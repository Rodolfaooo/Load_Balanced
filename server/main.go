package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

//Abre o server

func openServer() {

	s := &http.Server{
		Addr:           ":8080",
		Handler:        FunctionHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

//Cabeçalho da requisição

type FunctionHandler struct{}

func (f FunctionHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if !ValidServer(req) {
		res.Write(MessageToJson(TableMessage[404]))
		return
	}

	//metodo
	Method[req.Method][req.URL.Path](res, req)
}

func ValidServer(req *http.Request) bool {
	return Method[req.Method][req.URL.Path] != nil
}

type HTTPMessage struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

var TableMessage = map[int]HTTPMessage{
	404: {404, "not found"},
	500: {500, "Internal server error"},
}

func MessageToJson(m HTTPMessage) []byte {
	json, _ := json.Marshal(m)
	return json
}