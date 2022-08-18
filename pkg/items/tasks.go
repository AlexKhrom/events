package items

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
)

type TasksRepoInterface interface {
	NewTask(task *Task) bool
	GetTasks(userId int) ([]Task, bool)
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

type TaskPiece struct {
	Id     int     `json:"id"`
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
		fmt.Println("err new event sql2 = ", err)
		return false

	}

	lastID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("err new event sql3 = ", err)
		return false
	}

	fmt.Println("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	return false
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

func (r *TaskRepo) ChangeTask(task Task, userId int) error {

}

func (r *TaskRepo) MakeTaskPiece(taskPiece TaskPiece) bool {

}


