package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type student struct {
	Name  string `json:"name"`
	Class string `json:"class"`
}

var students = []student{}

func main() {
	http.HandleFunc("/students", studentsHandler)
	http.ListenAndServe(":8080", nil)
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createStudent(w, r)
	case "GET":
		getStudents(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newStudent student
	err := decoder.Decode(&newStudent)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	students = append(students, newStudent)
	w.WriteHeader(http.StatusCreated)
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(students)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", jsonData)
}
