package models

import "time"

type (
	SSScoresResponse struct {
		PlayerScores PlayerScores `json:"playerScores"`
		Metadata     struct {
			Total        int `json:"total"`
			Page         int `json:"page"`
			ItemsPerPage int `json:"itemsPerPage"`
		} `json:"metadata"`
	}
	PlayerScores []struct {
		Score struct {
			Id                    int         `json:"id"`
			LeaderboardPlayerInfo interface{} `json:"leaderboardPlayerInfo"`
			Rank                  int         `json:"rank"`
			BaseScore             int         `json:"baseScore"`
			ModifiedScore         int         `json:"modifiedScore"`
			Pp                    float64     `json:"pp"`
			Weight                float64     `json:"weight"`
			Modifiers             string      `json:"modifiers"`
			Multiplier            float64     `json:"multiplier"`
			BadCuts               int         `json:"badCuts"`
			MissedNotes           int         `json:"missedNotes"`
			MaxCombo              int         `json:"maxCombo"`
			FullCombo             bool        `json:"fullCombo"`
			Hmd                   int         `json:"hmd"`
			TimeSet               time.Time   `json:"timeSet"`
			HasReplay             bool        `json:"hasReplay"`
			DeviceHmd             string      `json:"deviceHmd"`
			DeviceControllerLeft  string      `json:"deviceControllerLeft"`
			DeviceControllerRight string      `json:"deviceControllerRight"`
		} `json:"score"`
		Leaderboard struct {
			Id              int    `json:"id"`
			SongHash        string `json:"songHash"`
			SongName        string `json:"songName"`
			SongSubName     string `json:"songSubName"`
			SongAuthorName  string `json:"songAuthorName"`
			LevelAuthorName string `json:"levelAuthorName"`
			Difficulty      struct {
				LeaderboardId int    `json:"leaderboardId"`
				Difficulty    int    `json:"difficulty"`
				GameMode      string `json:"gameMode"`
				DifficultyRaw string `json:"difficultyRaw"`
			} `json:"difficulty"`
			MaxScore          int         `json:"maxScore"`
			CreatedDate       time.Time   `json:"createdDate"`
			RankedDate        *time.Time  `json:"rankedDate"`
			QualifiedDate     *time.Time  `json:"qualifiedDate"`
			LovedDate         interface{} `json:"lovedDate"`
			Ranked            bool        `json:"ranked"`
			Qualified         bool        `json:"qualified"`
			Loved             bool        `json:"loved"`
			MaxPP             int         `json:"maxPP"`
			Stars             float64     `json:"stars"`
			Plays             int         `json:"plays"`
			DailyPlays        int         `json:"dailyPlays"`
			PositiveModifiers bool        `json:"positiveModifiers"`
			PlayerScore       interface{} `json:"playerScore"`
			CoverImage        string      `json:"coverImage"`
			Difficulties      interface{} `json:"difficulties"`
		} `json:"leaderboard"`
	}
	SSPlayersResponse struct {
		Players  []SSPlayerResponse `json:"players"`
		Metadata struct {
			Total        int `json:"total"`
			Page         int `json:"page"`
			ItemsPerPage int `json:"itemsPerPage"`
		} `json:"metadata"`
	}
	SSPlayerResponse struct {
		Id             string      `json:"id"`
		Name           string      `json:"name"`
		ProfilePicture string      `json:"profilePicture"`
		Bio            *string     `json:"bio"`
		Country        string      `json:"country"`
		Pp             float64     `json:"pp"`
		Rank           int         `json:"rank"`
		CountryRank    int         `json:"countryRank"`
		Role           interface{} `json:"role"`
		Badges         interface{} `json:"badges"`
		Histories      string      `json:"histories"`
		Permissions    int         `json:"permissions"`
		Banned         bool        `json:"banned"`
		Inactive       bool        `json:"inactive"`
		ScoreStats     struct {
			TotalScore            int64   `json:"totalScore"`
			TotalRankedScore      int     `json:"totalRankedScore"`
			AverageRankedAccuracy float64 `json:"averageRankedAccuracy"`
			TotalPlayCount        int     `json:"totalPlayCount"`
			RankedPlayCount       int     `json:"rankedPlayCount"`
			ReplaysWatched        int     `json:"replaysWatched"`
		} `json:"scoreStats"`
		FirstSeen time.Time `json:"firstSeen"`
	}
)
