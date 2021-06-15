package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Record struct {
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func main() {
	http.HandleFunc("/", createRecord)
	http.HandleFunc("/record", createRecord)
	http.HandleFunc("/records", getRecords)

	fmt.Println("server listening on http://localhost:80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Fprintf(os.Stderr, "server closed: %v", err)
		os.Exit(1)
	}
}

func createRecord(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		sendServerError(writer, err)
	}

	file, err := os.ReadFile("data/records.json")
	if err != nil {
		sendServerError(writer, err)
	}

	var records []Record
	if err := json.Unmarshal(file, &records); err != nil {
		sendServerError(writer, err)
	}

	records = append(records, Record{Headers: request.Header, Body: string(body)})
	recordsJSON, err := json.Marshal(records)
	if err != nil {
		sendServerError(writer, err)
	}

	if err := ioutil.WriteFile("data/records.json", recordsJSON, 0644); err != nil {
		sendServerError(writer, err)
	}

	writer.WriteHeader(200)
	writer.Header().Add("Content-Type", "text/plain")
	writer.Write([]byte("OK"))
}

func getRecords(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<h1>Logger</h1>"))
}

func sendServerError(writer http.ResponseWriter, err error) {
	writer.WriteHeader(500)
	writer.Header().Add("Content-Type", "text/plain")
	writer.Write([]byte(err.Error()))
}
