package models

type (
	Playlist struct {
		PlaylistTitle  string `json:"playlistTitle"`
		PlaylistAuthor string `json:"playlistAuthor"`
		CustomData     struct {
			SyncUrl string `json:"syncUrl"`
		} `json:"customData"`
		Image string              `json:"image"`
		Songs []PlaylistSongEntry `json:"songs"`
		Stats PlaylistStatsEntry  `json:"stats"`
	}
	PlaylistStatsEntry struct {
		SelfName         string `json:"selfName"`
		TargetName       string `json:"targetName"`
		SelfConsidered   int    `json:"selfConsidered"`
		TargetConsidered int    `json:"targetConsidered"`
		SnipeCount       int    `json:"snipeCount"`
		HoldCount        int    `json:"holdCount"`
	}
	PlaylistSongEntry struct {
		Hash         string                        `json:"hash"`
		Difficulties []PlaylistSongEntryDifficulty `json:"difficulties"`
		PpDiff       float64                       `json:"ppDifference"`
	}
	PlaylistSongEntryDifficulty struct {
		Name           string `json:"name"`
		Characteristic string `json:"characteristic"`
	}
)
