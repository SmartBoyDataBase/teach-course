package handler

import (
	"encoding/json"
	"net/http"
	"sbdb-teach-course/infrastructure"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	teacherId := r.URL.Query().Get("teacher_id")
	semesterId := r.URL.Query().Get("semester_id")
	rows, _ := infrastructure.DB.Query(`
	SELECT course_id from teachcourse
	WHERE teacher_id=$1 AND semester_id=$2;
	`, teacherId, semesterId)
	var result []uint64
	for rows.Next() {
		var courseId uint64
		rows.Scan(&courseId)
		result = append(result, courseId)
	}
	var body []byte
	if len(result) != 0 {
		body, _ = json.Marshal(result)
	} else {
		body = []byte("[]")
	}
	w.Write(body)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	teacherId := r.URL.Query().Get("teacher_id")
	semesterId := r.URL.Query().Get("semester_id")
	courseId := r.URL.Query().Get("course_id")
	_, err := infrastructure.DB.Exec(`
	INSERT INTO teachcourse(course_id, teacher_id, semester_id) VALUES ($1,$2,$3);
	`, courseId, teacherId, semesterId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		result := struct {
			TeacherId  string `json:"teacher_id"`
			SemesterId string `json:"semester_id"`
			CourseId   string `json:"course_id"`
		}{TeacherId: teacherId, SemesterId: semesterId, CourseId: courseId}
		body, _ := json.Marshal(result)
		_, _ = w.Write(body)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHandler(w, r)
	case "POST":
		postHandler(w, r)
	}
}

func AllHandler(w http.ResponseWriter, r *http.Request) {
	rows, _ := infrastructure.DB.Query(`
	SELECT course_id, teacher_id, semester_id
	FROM teachcourse;
	`)
	var result []struct {
		CourseId   uint64 `json:"course_id"`
		TeacherId  uint64 `json:"teacher_id"`
		SemesterId uint64 `json:"semester_id"`
	}
	for rows.Next() {
		var courseId uint64
		var teacherId uint64
		var semesterId uint64
		_ = rows.Scan(&courseId, &teacherId, &semesterId)
		result = append(result, struct {
			CourseId   uint64 `json:"course_id"`
			TeacherId  uint64 `json:"teacher_id"`
			SemesterId uint64 `json:"semester_id"`
		}{CourseId: courseId, TeacherId: teacherId, SemesterId: semesterId})
	}
	body, _ := json.Marshal(result)
	w.Write(body)
}
