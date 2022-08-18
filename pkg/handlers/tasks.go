package handlers

import (
	"database/sql"
	"encoding/json"
	"events/pkg/items"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TaskHandler struct {
	Repo items.TasksRepoInterface
}

func NewTaskHandler(db *sql.DB) *TaskHandler {
	handler := new(TaskHandler)
	handler.Repo = items.NewTaskRepo(db)
	return handler
}

func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi add task")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
		return
	}

	fmt.Println("request  = ", string(body))
	task := &items.Task{}

	err1 := json.Unmarshal(body, task)
	if err1 != nil {
		fmt.Println("some bad - ", err1)
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	userId, err3 := GetUserId(w, r)

	if !err3 {
		fmt.Println("some bad add task")
		return
	}

	task.UserId = userId

	fmt.Println("task = ", task)

	h.Repo.NewTask(task)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi get tasks")

	userId, err3 := GetUserId(w, r)

	if !err3 {
		fmt.Println("some bad add task")
		return
	}

	tasks, err := h.Repo.GetTasks(userId)
	if err {
		JSONError(w, http.StatusBadRequest, "some bad get tasks")
		return
	}

	fmt.Println("events = ", tasks)

	eventJson, err1 := json.Marshal(tasks)
	if err1 != nil {
		JSONError(w, http.StatusBadRequest, "some bad get tasks json.marshal")
		return
	}
	_, err1 = w.Write(eventJson)
	if err1 != nil {
		JSONError(w, http.StatusBadRequest, "some bad get tasks write response")
		return
	}

}

func (h *TaskHandler) ChangeTask(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi change event")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
	}
	fmt.Println("request  = ", string(body))
	event := &items.Event{}

	err1 := json.Unmarshal(body, event)
	if err1 != nil {
		fmt.Println("some bad")
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	fmt.Println("eventId = ", event)

	userId, ok := GetUserId(w, r)
	if !ok {
		fmt.Println("some bad get user ID")
		JSONError(w, http.StatusBadRequest, "some bad get user ID")
		return
	}

	fmt.Println("userId = ", userId)

	err = h.Repo.ChangeTask(event, userId)

	if err != nil {
		fmt.Println("some bad change events sql, err = ", err)
		JSONError(w, http.StatusBadRequest, "some bad delete events sql")
		return
	}

}
