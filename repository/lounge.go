package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type loungeRepository struct{}

func NewLoungeRepository() (LoungeRepository, error) {
	return loungeRepository{}, nil
}

func (r loungeRepository) GetLoungeName(ctx context.Context, userID string) (*GetLoungeNameResponse, error) {
	res, err := http.Get(r.getURLPlayerFromDiscord(userID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response GetLoungeNameResponse
	json.Unmarshal(body, &response)
	return &response, nil
}

func (r loungeRepository) getAPIURLBase() string {
	return "https://www.mk8dx-lounge.com/api"
}

func (r loungeRepository) getURLPlayerFromDiscord(discordUserID string) string {
	return r.getAPIURLBase() + "/player?discordId=" + discordUserID
}
