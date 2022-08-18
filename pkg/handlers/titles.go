package handlers

import (
	"database/sql"
	"encoding/json"
	"events/pkg/items"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TitlesHand struct {
	Repo items.TitlesRepoInterface
}

func NewTitlesHand(db *sql.DB) *TitlesHand {
	us := new(TitlesHand)
	us.Repo = items.NewTitleRepo(db)
	return us
}

func (h *TitlesHand) AddTitle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi add title ttt")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
	}
	fmt.Println("request  = ", string(body))
	title := &items.Title{}

	err1 := json.Unmarshal(body, title)
	if err1 != nil {
		fmt.Println("some bad")
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

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

	title.UserId = int(userId.(float64))
	fmt.Println("event = ", title)

	h.Repo.NewTitle(title)

}

func (h *TitlesHand) GetTitles(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi get titles")

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

	titles, err := h.Repo.GetTitles(int(userId.(float64)))
	if err {
		JSONError(w, http.StatusBadRequest, "some bad get events")
		return
	}

	fmt.Println("titles = ", titles)

	eventJson, err1 := json.Marshal(titles)
	if err1 != nil {
		JSONError(w, http.StatusBadRequest, "some bad get events json.marshal")
		return
	}
	_, err1 = w.Write(eventJson)
	if err1 != nil {
		JSONError(w, http.StatusBadRequest, "some bad get events write response")
		return
	}

}

func (h *TitlesHand) DeleteTitle(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi delete title")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
		return
	}
	fmt.Println("request  = ", string(body))
	titleId := &items.EventId{}

	err1 := json.Unmarshal(body, titleId)
	if err1 != nil {
		fmt.Println("some bad")
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	fmt.Println("titleID = ", titleId.Id)

	userId, ok := GetUserId(w, r)
	if !ok {
		JSONError(w, http.StatusBadRequest, "bad access token")
		return
	}

	fmt.Println("userId = ", userId)

	err = h.Repo.DeleteTitle(titleId.Id, userId)

	if err != nil {
		fmt.Println("some bad delete titles sql, err = ", err)
		JSONError(w, http.StatusBadRequest, "some bad delete events sql")
		return
	}

}
