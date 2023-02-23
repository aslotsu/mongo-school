package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

var _ = godotenv.Load()
var uri = os.Getenv("MONGODB_URI")
var serverApiOptions = options.ServerAPI(options.ServerAPIVersion1)
var clientOptions = options.Client().ApplyURI(uri).SetServerAPIOptions(serverApiOptions)
var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Hour)
var client, errorMsg = mongo.Connect(ctx, clientOptions)

type Course struct {
	//ID         primitive.ObjectID `bson:"id, omitempty" json:"id"`
	Name       string `bson:"name" json:"name"`
	CourseCode string `bson:"course_code" json:"course-code"`
	Class      int    `bson:"class" json:"class"`
	Pass       int    `bson:"pass" json:"pass"`
	Fail       int    `bson:"fail" json:"fail"`
}

func main() {
	defer cancel()
	if errorMsg != nil {
		log.Fatal(errorMsg)
	}
	_ = client.NumberSessionsInProgress()
	router := mux.NewRouter()
	router.HandleFunc("/new", newCourse).Methods("POST")
	router.HandleFunc("/new-ones", newCourses).Methods("POST")
	router.HandleFunc("/courses", getCourses).Methods("GET")
	router.HandleFunc("/course/:id", getCourse).Methods("GET")
	router.HandleFunc("/course/del/:id", removeCourse).Methods("DELETE")
	router.HandleFunc("/course/del", removeCourses).Methods("DELETE")
	log.Println("Your weird server is running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}

//
//func hello(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	fmt.Println("Hello")
//}

func newCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courses := client.Database("school").Collection("courses")
	var course Course
	err := json.NewDecoder(r.Body).Decode(&course)
	result, err := courses.InsertOne(ctx, course)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(w).Encode(course); err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.InsertedID)

}

func newCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = r
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = r
	coursesColl := client.Database("school").Collection("courses")
	var courses []Course
	cursor, err := coursesColl.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err := cursor.All(ctx, &courses); err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(courses)

}

func getCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application.json")
	_ = r
	name := mux.Vars(r)["name"]
	var course Course
	courseColl := client.Database("school").Collection("courses")
	_ = courseColl.FindOne(ctx, Course{Name: name})
	if err := json.NewEncoder(w).Encode(&course); err != nil {
		log.Fatal(err)
	}
}

func removeCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = r

}
func removeCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = r

}
