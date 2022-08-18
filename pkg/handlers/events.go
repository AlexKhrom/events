package handlers

import (
	"database/sql"
	"encoding/json"
	"events/pkg/items"
	"fmt"
	"io/ioutil"
	"net/http"
)

var i int

type EventsHand struct {
	Repo items.EventsRepoInterface
}

func NewEventsHand(db *sql.DB) *EventsHand {
	us := new(EventsHand)
	us.Repo = items.NewEventsRepo(db)
	return us
}

func (h *EventsHand) AddEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi add event ttt")

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

	event.UserId = int(userId.(float64))
	fmt.Println("event = ", event)

	h.Repo.NewEvent(event)

}

func (h *EventsHand) GetEvents(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi get events")

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

	events, err := h.Repo.GetEvents(int(userId.(float64)))
	if err {
		JSONError(w, http.StatusBadRequest, "some bad get events")
		return
	}

	fmt.Println("events = ", events)

	eventJson, err1 := json.Marshal(events)
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

func (h *EventsHand) ChangeEvent(w http.ResponseWriter, r *http.Request) {

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

	err = h.Repo.ChangeEvent(event, userId)

	if err != nil {
		fmt.Println("some bad change events sql, err = ", err)
		JSONError(w, http.StatusBadRequest, "some bad delete events sql")
		return
	}

}

func (h *EventsHand) DeleteEvent(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi delete event")

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

	err = h.Repo.DeleteEvent(eventId.Id, int(userId.(float64)))

	if err != nil {
		fmt.Println("some bad delete events sql, err = ", err)
		JSONError(w, http.StatusBadRequest, "some bad delete events sql")
		return
	}

	//fmt.Println("events = ", events)
	//
	//eventJson, err1 := json.Marshal(events)
	//if err1 != nil {
	//	JSONError(w, http.StatusBadRequest, "some bad get events json.marshal")
	//	return
	//}
	//_, err1 = w.Write(eventJson)
	//if err1 != nil {
	//	JSONError(w, http.StatusBadRequest, "some bad get events write response")
	//	return
	//}

}

func JSONError(w http.ResponseWriter, status int, msg string) {
	resp, err := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})
	w.WriteHeader(status)
	if err != nil {
		fmt.Println("error in JSONError ")
		return
	}
	_, err2 := w.Write(resp)
	if err2 != nil {
		fmt.Println("some bad in JSONError write response")
	}
}
