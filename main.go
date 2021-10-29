package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Model for course

type Course struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

// fake db

var courses []Course

func (c *Course) IsEmpty() bool {
	return c.CourseName == ""
}

func main() {
	r := mux.NewRouter()

	courses = append(courses, Course{"2", "React", 300, &Author{"charan", "charan.dev"}})

	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createCourse).Methods("POST")
	r.HandleFunc("/course/${id}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/course/${id}", deleteCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":4000", r))
}

// controller

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Golang API</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all Courses")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one course")
	w.Header().Set("Content-Type", "application/json")

	// get id from r
	params := mux.Vars(r)

	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	json.NewEncoder(w).Encode("No course found")

}

func createCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("Invalid data")
		return
	}

	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)

	if course.IsEmpty() {
		json.NewEncoder(w).Encode("No data")
		return
	}

	rand.Seed(time.Now().UnixNano())

	course.CourseId = strconv.Itoa(rand.Intn(100))

	courses = append(courses, course)

	json.NewEncoder(w).Encode(course)

}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	json.NewEncoder(w).Encode("No course found")
}

func deleteCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, course := range courses {
		if course.CourseId == params["Id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode(courses)
			break
		}
	}

	json.NewEncoder(w).Encode("No course found")
}
