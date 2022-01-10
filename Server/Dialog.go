package Server

import (
	"errors"
	"fmt"
	"net/http"
)

type DialogResponse struct {
	Type       string            `json:"type"`
	CallbackId string            `json:"callback_id"`
	State      string            `json:"state"`
	UserId     string            `json:"user_id"`
	ChannelId  string            `json:"channel_id"`
	TeamId     string            `json:"team_id"`
	Submission map[string]string `json:"submission"`
	Cancelled  bool              `json:"cancelled"`
}

func ParseDialogResponse(w http.ResponseWriter, req *http.Request) (DialogResponse, error) {
	var dialogResponse DialogResponse

	err := DecodeJSONBody(w, req, &dialogResponse)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return DialogResponse{}, err
	}

	if dialogResponse.Type != "dialog_submission" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return DialogResponse{}, fmt.Errorf("invalid request received: %s\n", dialogResponse.Type)
	}
	return dialogResponse, nil
}
