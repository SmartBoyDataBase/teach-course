package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sbdb-teach-course/infrastructure"
	"strconv"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := infrastructure.DB.QueryRow(`
	SELECT course_id, teacher_id, semester_id from teachcourse
	WHERE id=$1;
	`, id)
	var result struct {
		Id         uint64 `json:"id"`
		TeacherId  string `json:"teacher_id"`
		SemesterId string `json:"semester_id"`
		CourseId   string `json:"course_id"`
	}
	result.Id, _ = strconv.ParseUint(id, 10, 64)
	row.Scan(&result.CourseId, &result.TeacherId, &result.SemesterId)
	body, _ := json.Marshal(result)
	w.Write(body)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	teacherId := r.URL.Query().Get("teacher_id")
	semesterId := r.URL.Query().Get("semester_id")
	courseId := r.URL.Query().Get("course_id")
	row := infrastructure.DB.QueryRow(`
	INSERT INTO teachcourse(course_id, teacher_id, semester_id) 
	VALUES ($1,$2,$3)
	RETURNING id;`, courseId, teacherId, semesterId)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		result := struct {
			Id         uint64 `json:"id"`
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
	var rows *sql.Rows
	if r.URL.Query().Get("teacher_id") == "" {
		rows, _ = infrastructure.DB.Query(`
			SELECT id, course_id, teacher_id, semester_id
			FROM teachcourse;
		`)
	} else {
		rows, _ = infrastructure.DB.Query(`
			SELECT id, course_id, teacher_id, semester_id
			FROM teachcourse
			WHERE teacher_id=$1;
		`, r.URL.Query().Get("teacher_id"))
	}
	var result []struct {
		Id         uint64 `json:"id"`
		CourseId   uint64 `json:"course_id"`
		TeacherId  uint64 `json:"teacher_id"`
		SemesterId uint64 `json:"semester_id"`
	}
	for rows.Next() {
		var id uint64
		var courseId uint64
		var teacherId uint64
		var semesterId uint64
		_ = rows.Scan(&id, &courseId, &teacherId, &semesterId)
		result = append(result, struct {
			Id         uint64 `json:"id"`
			CourseId   uint64 `json:"course_id"`
			TeacherId  uint64 `json:"teacher_id"`
			SemesterId uint64 `json:"semester_id"`
		}{Id: id, CourseId: courseId, TeacherId: teacherId, SemesterId: semesterId})
	}
	body, _ := json.Marshal(result)
	w.Write(body)
}
