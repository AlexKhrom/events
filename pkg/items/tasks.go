package items

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type TasksRepoInterface interface {
	NewTask(task *Task) bool
	GetTasks(userId int) ([]Task, bool)
	GetTasksPiece(taskId, userId int) ([]TaskPiece, bool) //  return (tasksPiece , ok)
	ChangeTask(task *TaskForm, userId int) error
	MakeTaskPiece(task *TaskForm, userId int) bool
	DeleteTask(taskID, userId int) error
	DeleteTaskPiece(taskPieceID, userId int) error
}

type TaskRepo struct {
	DB  *sql.DB
	Mut sync.Mutex
}

type Task struct {
	Id        int     `json:"id"`
	UserId    int     `json:"userID"`
	Title     string  `json:"title"`
	TimeStart int     `json:"timeStart"`
	TimeEnd   int     `json:"timeEnd"`
	Type      string  `json:"type"`
	Target    float64 `json:"target"` //цель таска
	Num       float64 `json:"num"`    //колчиство выполненной части таска
}

type TaskForm struct {
	Id        int     `json:"id"`
	UserId    int     `json:"userID"`
	Title     string  `json:"title"`
	TimeStart int     `json:"timeStart"`
	TimeEnd   int     `json:"timeEnd"`
	Type      string  `json:"type"`
	Target    float64 `json:"target"` //цель таска
	Num       float64 `json:"num"`    //колчиство выполненной части таска
	DeltaNum  float64 `json:"deltaNum"`
}

type TaskPiece struct {
	Id     int     `json:"id"`
	UserId int     `json:"userId"`
	TaskId int     `json:"taskId"`
	DT     int     `json:"dt"`
	Type   string  `json:"type"`
	Num    float64 `json:"num"` //колчиство выполненной части таска
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	repo := new(TaskRepo)
	repo.DB = db
	return repo
}

func (r *TaskRepo) NewTask(task *Task) bool {

	result, err := r.DB.Exec(
		"INSERT INTO tasks (`user_id`,`title`,`time_start`,`time_end`,`type`,`target`,`num`) VALUES (?,?,?,?,?,?,?)",
		task.UserId,
		task.Title,
		task.TimeStart,
		task.TimeEnd,
		task.Type,
		task.Target,
		task.Num,
	)
	if err != nil {
		fmt.Println("err new task sql1 = ", err)
		return false
	}

	affected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("err new task sql2 = ", err)
		return false

	}

	lastID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("err new task sql3 = ", err)
		return false
	}

	fmt.Println("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	return true
}

func (r *TaskRepo) GetTasks(userId int) ([]Task, bool) {
	rows, err := r.DB.Query("SELECT * FROM tasks WHERE user_id =" + strconv.Itoa(userId))
	if err != nil {
		fmt.Println("get tasks error = ", err)
		return nil, true
	}

	var tasks []Task
	var i = 0

	for rows.Next() {

		var task = Task{}
		err = rows.Scan(&task.Id, &task.UserId, &task.Title, &task.TimeStart, &task.TimeEnd, &task.Type, &task.Target, &task.Num)
		if err != nil {
			fmt.Println("get tasks error = ", err)
			return nil, true
		}
		tasks = append(tasks, task)
		i++
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()
	return tasks, false
}

func (r *TaskRepo) GetTasksPiece(taskId, userId int) ([]TaskPiece, bool) {
	rows, err := r.DB.Query("SELECT * FROM taskPiece WHERE task_id =" + strconv.Itoa(taskId) + " AND user_id=" + strconv.Itoa(userId))
	if err != nil {
		fmt.Println("get tasks error = ", err)
		return nil, false
	}

	var tasksPieces []TaskPiece
	var i = 0

	for rows.Next() {

		var task = TaskPiece{}
		err = rows.Scan(&task.Id, &task.UserId, &task.TaskId, &task.DT, &task.Type, &task.Num)
		if err != nil {
			fmt.Println("get tasksPiece error = ", err)
			return nil, false
		}
		tasksPieces = append(tasksPieces, task)
		i++
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()
	return tasksPieces, true
}

func (r *TaskRepo) ChangeTask(task *TaskForm, userId int) error {
	fmt.Println("hi change task !!!")

	_, err := r.DB.Query("UPDATE tasks SET " + " time_start = " +
		strconv.Itoa(task.TimeStart) + ", time_end = " + strconv.Itoa(task.TimeEnd) + ", num = " + fmt.Sprintf("%f", task.Num) +
		" WHERE user_id = " + strconv.Itoa(userId) + " AND id = " + strconv.Itoa(task.Id))
	return err
}

func (r *TaskRepo) MakeTaskPiece(task *TaskForm, userId int) bool {
	result, err := r.DB.Exec(
		"INSERT INTO taskPiece (`user_id`,`task_id`,`dt`,`type`,`num`) VALUES (?,?,?,?,?)",
		userId,
		task.Id,
		time.Now().Unix(),
		task.Type,
		task.DeltaNum,
	)
	if err != nil {
		fmt.Println("err new taskPiece sql1 = ", err)
		return false
	}

	affected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("err new taskPiece sql2 = ", err)
		return false

	}

	lastID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("err new taskPiece sql3 = ", err)
		return false
	}

	fmt.Println("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	return true
}

func (r *TaskRepo) DeleteTask(taskID, userId int) error {
	_, err := r.DB.Query("DELETE FROM tasks WHERE user_id = " + strconv.Itoa(userId) + " AND id = " + strconv.Itoa(taskID))
	return err
}

func (r *TaskRepo) DeleteTaskPiece(taskPieceID, userId int) error {
	_, err := r.DB.Query("DELETE FROM taskPiece WHERE user_id = " + strconv.Itoa(userId) + " AND id = " + strconv.Itoa(taskPieceID))
	return err
}
