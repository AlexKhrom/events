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

func (h *TaskHandler) GetTasksPiece(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi get tasksPiece")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
	}
	fmt.Println("request  = ", string(body))
	taskId := &items.EventId{}

	err1 := json.Unmarshal(body, taskId)
	if err1 != nil {
		fmt.Println("some bad")
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	fmt.Println("taskId for pieces = ", taskId.Id)

	userId, err3 := GetUserId(w, r)

	if !err3 {
		fmt.Println("some bad add task")
		return
	}

	tasks, ok := h.Repo.GetTasksPiece(taskId.Id, userId)
	if !ok {
		JSONError(w, http.StatusBadRequest, "some bad get tasksPiece")
		return
	}

	fmt.Println("tasksPiece = ", tasks)

	eventJson, err1 := json.Marshal(tasks)
	if err1 != nil {
		JSONError(w, http.StatusBadRequest, "some bad get tasksPiece json.marshal")
		return
	}
	_, err1 = w.Write(eventJson)
	if err1 != nil {
		JSONError(w, http.StatusBadRequest, "some bad get tasksPiece write response")
		return
	}

}

func (h *TaskHandler) ChangeTask(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi change task")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
	}
	fmt.Println("request  = ", string(body))
	task := &items.TaskForm{}

	err1 := json.Unmarshal(body, task)
	if err1 != nil {
		fmt.Println("some bad")
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	fmt.Println("task  = ", task)

	userId, ok := GetUserId(w, r)
	if !ok {
		fmt.Println("some bad get user ID")
		JSONError(w, http.StatusBadRequest, "some bad get user ID")
		return
	}

	fmt.Println("userId = ", userId)

	err = h.Repo.ChangeTask(task, userId)

	if err != nil {
		fmt.Println("some bad change task sql, err = ", err)
		JSONError(w, http.StatusBadRequest, "some bad delete events sql")
		return
	}

	ok = h.Repo.MakeTaskPiece(task, userId)

	if !ok {
		fmt.Println("some bad change task new taskPiece sql, err = ", err)
		JSONError(w, http.StatusBadRequest, "some bad delete events sql")
		return
	}

}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi delete task")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
	}
	fmt.Println("request  = ", string(body))
	eventId := &items.EventId{}

	err1 := json.Unmarshal(body, eventId)
	if err1 != nil {
		fmt.Println("some bad")
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	fmt.Println("eventId = ", eventId.Id)

	tokenPaylod := items.GetTokenPayload(w, r)

	if tokenPaylod == nil {
		JSONError(w, http.StatusBadRequest, "bad access token")
		return
	}

	userId, ok := tokenPaylod["user"].(map[string]interface{})["id"]
	if !ok {
		JSONError(w, http.StatusBadRequest, "bad access token")
		return
	}

	fmt.Println("userId = ", userId)

	err = h.Repo.DeleteTask(eventId.Id, int(userId.(float64)))

	if err != nil {
		fmt.Println("some bad delete task sql, err = ", err)
		JSONError(w, http.StatusBadRequest, "some bad delete task sql")
		return
	}
}

func (h *TaskHandler) DeleteTaskPiece(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi delete taskPiece")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
	}
	fmt.Println("request  = ", string(body))
	eventId := &items.EventId{}

	err1 := json.Unmarshal(body, eventId)
	if err1 != nil {
		fmt.Println("some bad")
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	fmt.Println("eventId = ", eventId.Id)

	tokenPaylod := items.GetTokenPayload(w, r)

	if tokenPaylod == nil {
		JSONError(w, http.StatusBadRequest, "bad access token")
		return
	}

	userId, ok := tokenPaylod["user"].(map[string]interface{})["id"]
	if !ok {
		JSONError(w, http.StatusBadRequest, "bad access token")
		return
	}

	fmt.Println("userId = ", userId)

	err = h.Repo.DeleteTaskPiece(eventId.Id, int(userId.(float64)))

	if err != nil {
		fmt.Println("some bad delete taskPiece sql, err = ", err)
		JSONError(w, http.StatusBadRequest, "some bad delete taskPiece sql")
		return
	}
}
