package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type student struct {
	Name  string `json:"name"`
	Class string `json:"class"`
}

var students = []student{}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("studentDB").Collection("students")

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Swarup!")
	})

	router.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			createStudent(w, r, collection)
		case "GET":
			getStudents(w, r, collection)
		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	})

	router.HandleFunc("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		if r.Method != "DELETE" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

        fmt.Println(id)

		deleteStudent(w, r, collection, id)
	})

	fmt.Println("Server is listening on port 4000...")
	err = http.ListenAndServe(":4000", router)
	if err != nil {
		panic(err)
	}
}

func createStudent(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	decoder := json.NewDecoder(r.Body)
	var newStudent student // Change type from bson.M to student
	err := decoder.Decode(&newStudent)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//Insert new student into the collection
	insertResult, err := collection.InsertOne(context.Background(), newStudent)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println("New student added with ID:", insertResult.InsertedID)

	w.WriteHeader(http.StatusCreated)
}

func getStudents(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())

	var students []student

	// Iterate through the cursor
	for cur.Next(context.Background()) {
		var s student
		err := cur.Decode(&s)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}

	// Check for errors
	jsonData, err := json.Marshal(students)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON data to the response
	fmt.Fprintf(w, "%s", jsonData)
}

func deleteStudent(w http.ResponseWriter, r *http.Request, collection *mongo.Collection, id string) {
	// Convert the ID from a string to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Create a filter to find the student with the given ID
	filter := bson.M{"_id": objectID}

	// Delete the student from the collection
	deleteResult, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println("Student deleted with ID:", id)

	jsonData, err := json.Marshal(deleteResult)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON data to the response
	fmt.Fprintf(w, "%s", jsonData)
}
