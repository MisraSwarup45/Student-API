package main

import (
	"encoding/json"
	"fmt"
	"net/http"
    "log"
    "context"
    

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type student struct {
	Name  string `json:"name"`
	Class string `json:"class"`
}

var students = []student{}

func main() {

    clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    dbErr := client.Ping(context.Background(), nil)
    if dbErr != nil {
        log.Fatal(dbErr)
    }

    fmt.Println("Connected to MongoDB!")

    handler:= func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello Swarup!")
    }

    http.HandleFunc("/", handler);

	http.HandleFunc("/students", studentsHandler)
    fmt.Println("Server is listening on port 8080...")

	connectionErr := http.ListenAndServe(":8080", nil)
    if connectionErr != nil {
        panic(connectionErr)
    }
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
    fmt.Println("New student added: ", newStudent)
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
