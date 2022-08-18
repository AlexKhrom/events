package items

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
)

type TitlesRepoInterface interface {
	NewTitle(title *Title) bool
	GetTitles(userId int) ([]Title, bool)
	DeleteTitle(titleId, userId int) error
}

type TitleRepo struct {
	DB  *sql.DB
	Mut sync.Mutex
}

type Title struct {
	Id           int    `json:"id"`
	UserId       int    `json:"userId"`
	Title        string `json:"title"`
	WithDuration bool   `json:"withDuration"`
}

func NewTitleRepo(db *sql.DB) *TitleRepo {
	repo := new(TitleRepo)
	repo.DB = db
	return repo
}

func (r *TitleRepo) NewTitle(title *Title) bool {

	fmt.Println("title = ", title.Title)
	result, err := r.DB.Exec(
		"INSERT INTO titles (`user_id`,`title`,`with_duration`) VALUES (?,?,?)",
		title.UserId,
		title.Title,
		title.WithDuration,
	)
	if err != nil {
		fmt.Println("err new title sql1 = ", err)
		return false
	}

	affected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("err new title sql2 = ", err)
		return false

	}

	lastID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("err new title sql3 = ", err)
		return false
	}

	fmt.Println("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	return false
}

func (r *TitleRepo) GetTitles(userId int) ([]Title, bool) {
	fmt.Println("hi there!")
	rows, err := r.DB.Query("SELECT * FROM titles WHERE user_id =" + strconv.Itoa(userId))
	if err != nil {
		fmt.Println("get titles error = ", err)
		return nil, true
	}

	var titles []Title
	var i = 0

	for rows.Next() {

		var title = Title{}
		err = rows.Scan(&title.Id, &title.UserId, &title.Title, &title.WithDuration)
		if err != nil {
			fmt.Println("get titles error = ", err)
			return nil, true
		}
		titles = append(titles, title)
		i++
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()
	return titles, false
}

func (r *TitleRepo) DeleteTitle(titleId, userId int) error {
	_, err := r.DB.Query("DELETE FROM titles WHERE user_id = " + strconv.Itoa(userId) + " AND id = " + strconv.Itoa(titleId))
	return err
}
