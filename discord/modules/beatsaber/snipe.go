package beatsaber

import (
	"HelpingPixl/config"
	"HelpingPixl/models"
	"HelpingPixl/utils"
	"cmp"
	"fmt"
	"github.com/disgoorg/disgo/events"
	"log/slog"
	"math"
	"slices"
	"strconv"
)

func SnipeHoldPlaylist(e *events.ApplicationCommandInteractionCreate, self string, target string, leaderboard int) (models.Playlist, models.Playlist, string, bool) {
	var snipePlaylist models.Playlist
	var holdPlaylist models.Playlist
	var validSnipeMaps []models.PlaylistSongEntry
	var validHoldMaps []models.PlaylistSongEntry

	switch leaderboard {
	case 0:
		selfPlayerStructs, err := utils.FetchToStruct[models.BLPlayersResponse](fmt.Sprintf(models.BLBase+models.BLPlayers, self))
		if err != nil {
			slog.Error(err.Error())
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.InternalError), false
		}

		targetPlayerStructs, err := utils.FetchToStruct[models.BLPlayersResponse](fmt.Sprintf(models.BLBase+models.BLPlayers, target))
		if err != nil {
			slog.Error(err.Error())
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.InternalError), false
		}

		if !(len(selfPlayerStructs.Data) > 0) {
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.NoPlayerFound, self), false
		}
		if !(len(targetPlayerStructs.Data) > 0) {
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.NoPlayerFound, target), false
		}

		selfStruct := selfPlayerStructs.Data[0]
		targetStruct := targetPlayerStructs.Data[0]

		_ = e.DeferCreateMessage(true)

		//Fetching Scores

		//Self
		err, selfScoresStructs := collectBLScores(selfStruct.Id)
		if err != nil {
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.FetchingScoresFailed, selfStruct.Name), true
		}

		//Target
		err, targetScoresStructs := collectBLScores(targetStruct.Id)
		if err != nil {
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.FetchingScoresFailed, targetStruct.Name), true
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

		snipePlaylist.Stats = stats
		holdPlaylist.Stats = stats
		snipePlaylist.PlaylistTitle = fmt.Sprintf("Snipe BL (%s ▶ %s)", selfStruct.Name, targetStruct.Name)
		holdPlaylist.PlaylistTitle = fmt.Sprintf("Hold BL (%s ● %s)", selfStruct.Name, targetStruct.Name)
		break
	case 1:

		selfPlayerStructs, err := utils.FetchToStruct[models.SSPlayersResponse](fmt.Sprintf(models.SSBase+models.SSPlayers, self))
		if err != nil {
			slog.Error(err.Error())
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.InternalError), false
		}

		targetPlayerStructs, err := utils.FetchToStruct[models.SSPlayersResponse](fmt.Sprintf(models.SSBase+models.SSPlayers, target))
		if err != nil {
			slog.Error(err.Error())
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.InternalError), false
		}

		if !(len(selfPlayerStructs.Players) > 0) {
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.NoPlayerFound, self), false
		}
		if !(len(targetPlayerStructs.Players) > 0) {
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.NoPlayerFound, target), false
		}

		selfStruct := selfPlayerStructs.Players[0]
		targetStruct := targetPlayerStructs.Players[0]

		_ = e.DeferCreateMessage(true)

		//Self
		err, selfScoresStructs := collectSSScores(selfStruct.Id)
		if err != nil {
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.FetchingScoresFailed, selfStruct.Name), true
		}

		//Target
		err, targetScoresStructs := collectSSScores(targetStruct.Id)
		if err != nil {
			return models.Playlist{}, models.Playlist{}, fmt.Sprintf(config.Config.Formatting.FetchingScoresFailed, targetStruct.Name), true
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

		snipePlaylist.Stats = stats
		holdPlaylist.Stats = stats
		snipePlaylist.PlaylistTitle = fmt.Sprintf("Snipe SS (%s ▶ %s)", selfStruct.Name, targetStruct.Name)
		holdPlaylist.PlaylistTitle = fmt.Sprintf("Hold SS (%s ● %s)", selfStruct.Name, targetStruct.Name)
	}

	//Playlist Packing

	slices.SortFunc(validSnipeMaps, func(a, b models.PlaylistSongEntry) int {
		return cmp.Compare(b.PpDiff, a.PpDiff)
	})
	slices.SortFunc(validHoldMaps, func(a, b models.PlaylistSongEntry) int {
		return cmp.Compare(a.PpDiff, b.PpDiff)
	})

	snipePlaylist.PlaylistAuthor = "PixlPainter"
	snipePlaylist.Image = config.Config.BeatSaber.SnipeImage
	snipePlaylist.CustomData.SyncUrl = config.Config.BeatSaber.SnipeSyncUrl
	snipePlaylist.Songs = validSnipeMaps

	holdPlaylist.PlaylistAuthor = "PixlPainter"
	holdPlaylist.Image = config.Config.BeatSaber.HoldImage
	holdPlaylist.CustomData.SyncUrl = config.Config.BeatSaber.HoldSyncUrl
	holdPlaylist.Songs = validHoldMaps

	return snipePlaylist, holdPlaylist, "", true
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
