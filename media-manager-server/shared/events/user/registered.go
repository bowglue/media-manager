package user

import (
	"encoding/json"
	"shared/events"
)

type UserRegisterEvent struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func Publish(event UserRegisterEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return events.Publish("movies.created", data)
}
