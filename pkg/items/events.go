package items

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
)

type Event struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userID"`
	TitleId   int    `json:"titleId"`
	Title     string `json:"title"`
	DT        int    `json:"dt"`
	TimeStart int    `json:"timeStart"`
	TimeEnd   int    `json:"timeEnd"`
	Comment   string `json:"comment"`
}

type EventId struct {
	Id int `json:"id"`
}

type EventsRepoInterface interface {
	NewEvent(event *Event) bool
	GetEvents(userId int) ([]Event, bool)
	ChangeEvent(event *Event, userId int) error
	DeleteEvent(eventId, userId int) error
}

type EventsRepo struct {
	DB  *sql.DB
	Mut sync.Mutex
}

func NewEventsRepo(db *sql.DB) *EventsRepo {
	repo := new(EventsRepo)
	repo.DB = db
	return repo
}

func (r *EventsRepo) NewEvent(event *Event) bool {

	result, err := r.DB.Exec(
		"INSERT INTO events (`user_id`,`title_id`,`title`, `dt`,`time_start`,`time_end`) VALUES (?,?,?,?,?,?)",
		event.UserId,
		event.TitleId,
		event.Title,
		event.DT,
		event.TimeStart,
		event.TimeEnd,
	)
	if err != nil {
		fmt.Println("err new event sql1 = ", err)
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

func (r *EventsRepo) GetEvents(userId int) ([]Event, bool) {
	rows, err := r.DB.Query("SELECT * FROM events WHERE user_id =" + strconv.Itoa(userId))
	if err != nil {
		fmt.Println("get events error = ", err)
		return nil, true
	}

	var events []Event
	var i = 0

	for rows.Next() {

		var event = Event{}
		err = rows.Scan(&event.Id, &event.UserId, &event.TitleId, &event.Title, &event.DT, &event.TimeStart, &event.TimeEnd, &event.Comment)
		if err != nil {
			fmt.Println("get events error = ", err)
			return nil, true
		}
		events = append(events, event)
		i++
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()
	return events, false
}

func (r *EventsRepo) ChangeEvent(event *Event, userId int) error {
	fmt.Println("hi change event !!!")

	var comment string
	if event.Comment == "" {
		comment = "no comment"
	} else {
		comment = event.Comment
	}
	fmt.Println("event = ", event)
	_, err := r.DB.Query("UPDATE events SET dt = " + strconv.Itoa(event.DT) + ", time_start = " +
		strconv.Itoa(event.TimeStart) + ", time_end = " + strconv.Itoa(event.TimeEnd) + ", comment = '" + Ecran(comment) +
		"' WHERE user_id = " + strconv.Itoa(userId) + " AND id = " + strconv.Itoa(event.Id))
	return err
}

func (r *EventsRepo) DeleteEvent(eventId, userId int) error {
	_, err := r.DB.Query("DELETE FROM events WHERE user_id = " + strconv.Itoa(userId) + " AND id = " + strconv.Itoa(eventId))
	return err
}
