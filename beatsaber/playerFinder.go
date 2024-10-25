package beatsaber

import (
	"HelpingPixl/config"
	"HelpingPixl/models"
	"HelpingPixl/utils"
	"errors"
	"fmt"
	"log/slog"
)

func FindBLPlayerByName(tag string) (*models.BLPlayerResponse, error) {
	playerStructs, err := utils.FetchToStruct[models.BLPlayersResponse](fmt.Sprintf(models.BLBase+models.BLPlayersQuery, tag))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	if !(len(playerStructs.Data) > 0) {
		return nil, errors.New(config.Config.Formatting.NoPlayerFound)
	}
	return &playerStructs.Data[0], nil
}
func FindSSPlayerByName(tag string) (*models.SSPlayerResponse, error) {
	playerStructs, err := utils.FetchToStruct[models.SSPlayersResponse](fmt.Sprintf(models.SSBase+models.SSPlayers, tag))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	if !(len(playerStructs.Players) > 0) {
		return nil, errors.New(config.Config.Formatting.NoPlayerFound)
	}
	return &playerStructs.Players[0], nil
}
func GetBLPlayerById(id string) (*models.BLPlayerResponse, error) {
	playerStruct, err := utils.FetchToStruct[models.BLPlayerResponse](fmt.Sprintf(models.BLBase+models.BLPlayer, id))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &playerStruct, nil
}
func GetSSPlayerById(id string) (*models.SSPlayerResponse, error) {
	playerStruct, err := utils.FetchToStruct[models.SSPlayerResponse](fmt.Sprintf(models.SSBase+models.SSPlayer, id))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &playerStruct, nil
}
