package beatsaber

import (
	"HelpingPixl/config"
	"HelpingPixl/models"
	"HelpingPixl/utils"
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strconv"
)

func SnipeHoldPlaylist(selfTag, targetTag, selfId, targetId *string, leaderboard int) (snipePlaylist, holdPlaylist *models.Playlist, errorStr string) {
	var validSnipeMaps []models.PlaylistSongEntry
	var validHoldMaps []models.PlaylistSongEntry

	var tempSnipePlaylist models.Playlist
	var tempHoldPlaylist models.Playlist

	switch leaderboard {
	case 0:
		var selfStruct *models.BLPlayerResponse
		var targetStruct *models.BLPlayerResponse

		if selfId != nil && targetId != nil {
			selfStruct, _ = GetBLPlayerById(*selfId)
			targetStruct, _ = GetBLPlayerById(*targetId)
		} else if selfTag != nil && targetTag != nil {
			selfStruct, _ = FindBLPlayerByName(*selfTag)
			targetStruct, _ = FindBLPlayerByName(*targetTag)
		} else {
			return nil, nil, config.Config.Formatting.NoPlayerFound
		}

		if selfStruct == nil || targetStruct == nil {
			return nil, nil, config.Config.Formatting.NoPlayerFound
		}

		selfId = &selfStruct.Id
		targetId = &targetStruct.Id

		//Fetching Scores

		//Self
		err, selfScoresStructs := collectBLScores(selfStruct.Id)
		if err != nil {
			return nil, nil, fmt.Sprintf(config.Config.Formatting.FetchingScoresFailed, selfStruct.Name)
		}

		//Target
		err, targetScoresStructs := collectBLScores(targetStruct.Id)
		if err != nil {
			return nil, nil, fmt.Sprintf(config.Config.Formatting.FetchingScoresFailed, targetStruct.Name)
		}

		//Playlist Generation

		for _, sEntry := range selfScoresStructs.Data {
			for _, tEntry := range targetScoresStructs.Data {
				if sEntry.Leaderboard.Id == tEntry.Leaderboard.Id {
					dif := sEntry.Score.Pp - tEntry.Score.Pp

					if dif <= 0 {
						validSnipeMaps = append(validSnipeMaps, models.PlaylistSongEntry{
							Hash:         sEntry.Leaderboard.SongHash,
							Difficulties: []models.PlaylistSongEntryDifficulty{{Name: formatDifficulty(sEntry.Leaderboard.Difficulty), Characteristic: sEntry.Leaderboard.ModeName}},
							PpDiff:       dif,
						})
					} else {
						validHoldMaps = append(validHoldMaps, models.PlaylistSongEntry{
							Hash:         sEntry.Leaderboard.SongHash,
							Difficulties: []models.PlaylistSongEntryDifficulty{{Name: formatDifficulty(sEntry.Leaderboard.Difficulty), Characteristic: sEntry.Leaderboard.ModeName}},
							PpDiff:       dif,
						})
					}

				}
			}
		}

		stats := models.PlaylistStatsEntry{
			SelfName:         selfStruct.Name,
			TargetName:       targetStruct.Name,
			SelfConsidered:   len(selfScoresStructs.Data),
			TargetConsidered: len(targetScoresStructs.Data),
			SnipeCount:       len(validSnipeMaps),
			HoldCount:        len(validHoldMaps),
		}

		tempSnipePlaylist.Stats = stats
		tempHoldPlaylist.Stats = stats
		tempSnipePlaylist.PlaylistTitle = fmt.Sprintf("Snipe BL (%s ▶ %s)", selfStruct.Name, targetStruct.Name)
		tempHoldPlaylist.PlaylistTitle = fmt.Sprintf("Hold BL (%s ● %s)", selfStruct.Name, targetStruct.Name)
		break
	case 1:

		var selfStruct *models.SSPlayerResponse
		var targetStruct *models.SSPlayerResponse

		if selfId != nil && targetId != nil {
			selfStruct, _ = GetSSPlayerById(*selfId)
			targetStruct, _ = GetSSPlayerById(*targetId)
		} else if selfTag != nil && targetTag != nil {
			selfStruct, _ = FindSSPlayerByName(*selfTag)
			targetStruct, _ = FindSSPlayerByName(*targetTag)
		} else {
			return nil, nil, config.Config.Formatting.NoPlayerFound
		}

		if selfStruct == nil || targetStruct == nil {
			return nil, nil, config.Config.Formatting.NoPlayerFound
		}

		selfId = &selfStruct.Id
		targetId = &targetStruct.Id

		//Self
		err, selfScoresStructs := collectSSScores(selfStruct.Id)
		if err != nil {
			return nil, nil, fmt.Sprintf(config.Config.Formatting.FetchingScoresFailed, selfStruct.Name)
		}

		//Target
		err, targetScoresStructs := collectSSScores(targetStruct.Id)
		if err != nil {
			return nil, nil, fmt.Sprintf(config.Config.Formatting.FetchingScoresFailed, targetStruct.Name)
		}

		for _, sEntry := range selfScoresStructs.PlayerScores {
			for _, tEntry := range targetScoresStructs.PlayerScores {
				if sEntry.Leaderboard.Id == tEntry.Leaderboard.Id {
					dif := sEntry.Score.Pp - tEntry.Score.Pp
					if dif <= 0 {
						validSnipeMaps = append(validSnipeMaps, models.PlaylistSongEntry{
							Hash:         sEntry.Leaderboard.SongHash,
							Difficulties: []models.PlaylistSongEntryDifficulty{{Name: formatDifficulty(sEntry.Leaderboard.Difficulty.Difficulty), Characteristic: formatGameMode(sEntry.Leaderboard.Difficulty.GameMode)}},
							PpDiff:       dif,
						})
					} else {
						validHoldMaps = append(validHoldMaps, models.PlaylistSongEntry{
							Hash:         sEntry.Leaderboard.SongHash,
							Difficulties: []models.PlaylistSongEntryDifficulty{{Name: formatDifficulty(sEntry.Leaderboard.Difficulty.Difficulty), Characteristic: sEntry.Leaderboard.Difficulty.GameMode}},
							PpDiff:       dif,
						})
					}

				}
			}
		}

		stats := models.PlaylistStatsEntry{
			SelfName:         selfStruct.Name,
			TargetName:       targetStruct.Name,
			SelfConsidered:   len(selfScoresStructs.PlayerScores),
			TargetConsidered: len(targetScoresStructs.PlayerScores),
			SnipeCount:       len(validSnipeMaps),
			HoldCount:        len(validHoldMaps),
		}

		tempSnipePlaylist.Stats = stats
		tempHoldPlaylist.Stats = stats
		tempSnipePlaylist.PlaylistTitle = fmt.Sprintf("Snipe SS (%s ▶ %s)", selfStruct.Name, targetStruct.Name)
		tempHoldPlaylist.PlaylistTitle = fmt.Sprintf("Hold SS (%s ● %s)", selfStruct.Name, targetStruct.Name)
	}

	//Playlist Packing

	slices.SortFunc(validSnipeMaps, func(a, b models.PlaylistSongEntry) int {
		return cmp.Compare(b.PpDiff, a.PpDiff)
	})
	slices.SortFunc(validHoldMaps, func(a, b models.PlaylistSongEntry) int {
		return cmp.Compare(a.PpDiff, b.PpDiff)
	})

	tempSnipePlaylist.PlaylistAuthor = "PixlPainter"
	tempSnipePlaylist.Image = config.Config.BeatSaber.SnipeImage
	tempSnipePlaylist.CustomData.SyncUrl = buildSyncUrl(*selfId, *targetId, "snipe", leaderboard)
	tempSnipePlaylist.Songs = validSnipeMaps

	tempHoldPlaylist.PlaylistAuthor = "PixlPainter"
	tempHoldPlaylist.Image = config.Config.BeatSaber.HoldImage
	tempHoldPlaylist.CustomData.SyncUrl = buildSyncUrl(*selfId, *targetId, "hold", leaderboard)
	tempHoldPlaylist.Songs = validHoldMaps

	return &tempSnipePlaylist, &tempHoldPlaylist, ""
}

func buildSyncUrl(self, target, plType string, leaderboard int) string {
	if config.Config.WebServerAPI.APIUrl == "" {
		return ""
	}
	return fmt.Sprintf("%s/beatsaber/playlist/%s?self=%s&target=%s&leaderboard=%d", config.Config.WebServerAPI.APIUrl, plType, self, target, leaderboard)
}

func formatGameMode(mode string) string {
	switch mode {
	case "SoloStandard":
		return "Standard"
	}
	return "Standard"
}
func formatDifficulty(diff int) string {
	switch diff {
	case 1:
		return "easy"
	case 3:
		return "normal"
	case 5:
		return "hard"
	case 7:
		return "expert"
	case 9:
		return "expertPlus"
	}
	return ""
}

func collectBLScores(id string) (err error, scoresResponse models.BLScoresResponse) {
	scoresResponse, err = utils.FetchToStruct[models.BLScoresResponse](fmt.Sprintf(models.BLBase+models.BLPlayerScores, id, 1))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	var page = 2
	var totalPages = int(math.Floor(float64(scoresResponse.Metadata.Total/scoresResponse.Metadata.ItemsPerPage)) + 1)

	for page <= totalPages {
		slog.Debug("Query Iteration:" + strconv.Itoa(page))

		temp, err := utils.FetchToStruct[models.BLScoresResponse](fmt.Sprintf(models.BLBase+models.BLPlayerScores, id, page))
		if err != nil {
			slog.Error(err.Error())
			return err, scoresResponse
		}
		for _, entry := range temp.Data {
			scoresResponse.Data = append(scoresResponse.Data, entry)
		}
		page++
	}
	return
}
func collectSSScores(id string) (err error, scoresResponse models.SSScoresResponse) {
	scoresResponse, err = utils.FetchToStruct[models.SSScoresResponse](fmt.Sprintf(models.SSBase+models.SSPlayerScores, id))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if !(scoresResponse.Metadata.Total/scoresResponse.Metadata.ItemsPerPage > 1) {
		return
	}

	var page = 2
	var totalPages = int(math.Floor(float64(scoresResponse.Metadata.Total/scoresResponse.Metadata.ItemsPerPage)) + 1)
	var stillRanked bool = true

	for page <= totalPages {
		if !stillRanked {
			break
		}
		slog.Debug("Query Iteration:" + strconv.Itoa(page))

		temp, err := utils.FetchToStruct[models.SSScoresResponse](fmt.Sprintf(models.SSBase+models.SSPlayerScores+"&page=%d", id, page))
		if err != nil {
			slog.Error(err.Error())
			return err, scoresResponse
		}
		for _, entry := range temp.PlayerScores {
			if !entry.Leaderboard.Ranked {
				stillRanked = false
			} else {
				scoresResponse.PlayerScores = append(scoresResponse.PlayerScores, entry)
			}
		}

		page++
	}

	return
}
