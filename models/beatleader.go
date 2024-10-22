package models

type BLSongsResponse struct {
	Metadata struct {
		ItemsPerPage int `json:"itemsPerPage"`
		Page         int `json:"page"`
		Total        int `json:"total"`
	} `json:"metadata"`
	Data []struct {
		ValidContexts int `json:"validContexts"`
		Leaderboard   struct {
			Id   string `json:"id"`
			Song struct {
				Id              string      `json:"id"`
				Hash            string      `json:"hash"`
				Name            string      `json:"name"`
				SubName         string      `json:"subName"`
				Author          string      `json:"author"`
				Mapper          string      `json:"mapper"`
				MapperId        int         `json:"mapperId"`
				CollaboratorIds interface{} `json:"collaboratorIds"`
				CoverImage      string      `json:"coverImage"`
				Bpm             float64     `json:"bpm"`
				Duration        int         `json:"duration"`
				FullCoverImage  string      `json:"fullCoverImage"`
			} `json:"song"`
		} `json:"leaderboard"`
		AccLeft        float64     `json:"accLeft"`
		AccRight       float64     `json:"accRight"`
		Id             int         `json:"id"`
		BaseScore      int         `json:"baseScore"`
		ModifiedScore  int         `json:"modifiedScore"`
		Accuracy       float64     `json:"accuracy"`
		PlayerId       string      `json:"playerId"`
		Pp             float64     `json:"pp"`
		BonusPp        float64     `json:"bonusPp"`
		PassPP         float64     `json:"passPP"`
		AccPP          float64     `json:"accPP"`
		TechPP         float64     `json:"techPP"`
		Rank           int         `json:"rank"`
		Country        string      `json:"country"`
		FcAccuracy     float64     `json:"fcAccuracy"`
		FcPp           float64     `json:"fcPp"`
		Weight         float64     `json:"weight"`
		Replay         string      `json:"replay"`
		Modifiers      string      `json:"modifiers"`
		BadCuts        int         `json:"badCuts"`
		MissedNotes    int         `json:"missedNotes"`
		BombCuts       int         `json:"bombCuts"`
		WallsHit       int         `json:"wallsHit"`
		Pauses         int         `json:"pauses"`
		FullCombo      bool        `json:"fullCombo"`
		Platform       string      `json:"platform"`
		MaxCombo       int         `json:"maxCombo"`
		MaxStreak      *int        `json:"maxStreak"`
		Hmd            int         `json:"hmd"`
		Controller     int         `json:"controller"`
		LeaderboardId  string      `json:"leaderboardId"`
		Timeset        string      `json:"timeset"`
		Timepost       int         `json:"timepost"`
		ReplaysWatched int         `json:"replaysWatched"`
		PlayCount      int         `json:"playCount"`
		LastTryTime    int         `json:"lastTryTime"`
		Priority       int         `json:"priority"`
		Player         interface{} `json:"player"`
	} `json:"data"`
}
type BLPlayersResponse struct {
	Metadata struct {
		ItemsPerPage int `json:"itemsPerPage"`
		Page         int `json:"page"`
		Total        int `json:"total"`
	} `json:"metadata"`
	Data []BLPlayerResponse `json:"data"`
}
type BLPlayerResponse struct {
	AccPp           float64 `json:"accPp"`
	PassPp          float64 `json:"passPp"`
	TechPp          float64 `json:"techPp"`
	ProfileSettings struct {
		Id                    int         `json:"id"`
		Bio                   interface{} `json:"bio"`
		Message               interface{} `json:"message"`
		EffectName            string      `json:"effectName"`
		ProfileAppearance     string      `json:"profileAppearance"`
		Hue                   *float32    `json:"hue"`
		Saturation            *float32    `json:"saturation"`
		LeftSaberColor        interface{} `json:"leftSaberColor"`
		RightSaberColor       interface{} `json:"rightSaberColor"`
		ProfileCover          interface{} `json:"profileCover"`
		StarredFriends        string      `json:"starredFriends"`
		HorizontalRichBio     bool        `json:"horizontalRichBio"`
		RankedMapperSort      *string     `json:"rankedMapperSort"`
		ShowBots              bool        `json:"showBots"`
		ShowAllRatings        bool        `json:"showAllRatings"`
		ShowStatsPublic       bool        `json:"showStatsPublic"`
		ShowStatsPublicPinned bool        `json:"showStatsPublicPinned"`
	} `json:"profileSettings"`
	ScoreStats *struct {
		Id                            int     `json:"id"`
		TotalScore                    int     `json:"totalScore"`
		TotalUnrankedScore            int     `json:"totalUnrankedScore"`
		TotalRankedScore              int     `json:"totalRankedScore"`
		LastScoreTime                 int     `json:"lastScoreTime"`
		LastUnrankedScoreTime         int     `json:"lastUnrankedScoreTime"`
		LastRankedScoreTime           int     `json:"lastRankedScoreTime"`
		AverageRankedAccuracy         float64 `json:"averageRankedAccuracy"`
		AverageWeightedRankedAccuracy float64 `json:"averageWeightedRankedAccuracy"`
		AverageUnrankedAccuracy       float64 `json:"averageUnrankedAccuracy"`
		AverageAccuracy               float64 `json:"averageAccuracy"`
		MedianRankedAccuracy          float64 `json:"medianRankedAccuracy"`
		MedianAccuracy                float64 `json:"medianAccuracy"`
		TopRankedAccuracy             float64 `json:"topRankedAccuracy"`
		TopUnrankedAccuracy           float64 `json:"topUnrankedAccuracy"`
		TopAccuracy                   float64 `json:"topAccuracy"`
		TopPp                         float64 `json:"topPp"`
		TopBonusPP                    float64 `json:"topBonusPP"`
		TopPassPP                     float64 `json:"topPassPP"`
		TopAccPP                      float64 `json:"topAccPP"`
		TopTechPP                     float64 `json:"topTechPP"`
		PeakRank                      int     `json:"peakRank"`
		RankedMaxStreak               int     `json:"rankedMaxStreak"`
		UnrankedMaxStreak             int     `json:"unrankedMaxStreak"`
		MaxStreak                     int     `json:"maxStreak"`
		AverageLeftTiming             float64 `json:"averageLeftTiming"`
		AverageRightTiming            float64 `json:"averageRightTiming"`
		RankedPlayCount               int     `json:"rankedPlayCount"`
		UnrankedPlayCount             int     `json:"unrankedPlayCount"`
		TotalPlayCount                int     `json:"totalPlayCount"`
		RankedImprovementsCount       int     `json:"rankedImprovementsCount"`
		UnrankedImprovementsCount     int     `json:"unrankedImprovementsCount"`
		TotalImprovementsCount        int     `json:"totalImprovementsCount"`
		RankedTop1Count               int     `json:"rankedTop1Count"`
		UnrankedTop1Count             int     `json:"unrankedTop1Count"`
		Top1Count                     int     `json:"top1Count"`
		RankedTop1Score               int     `json:"rankedTop1Score"`
		UnrankedTop1Score             int     `json:"unrankedTop1Score"`
		Top1Score                     int     `json:"top1Score"`
		AverageRankedRank             float64 `json:"averageRankedRank"`
		AverageWeightedRankedRank     float64 `json:"averageWeightedRankedRank"`
		AverageUnrankedRank           float64 `json:"averageUnrankedRank"`
		AverageRank                   float64 `json:"averageRank"`
		SspPlays                      int     `json:"sspPlays"`
		SsPlays                       int     `json:"ssPlays"`
		SpPlays                       int     `json:"spPlays"`
		SPlays                        int     `json:"sPlays"`
		APlays                        int     `json:"aPlays"`
		TopPlatform                   string  `json:"topPlatform"`
		TopHMD                        int     `json:"topHMD"`
		AllHMDs                       string  `json:"allHMDs"`
		TopPercentile                 float64 `json:"topPercentile"`
		CountryTopPercentile          float64 `json:"countryTopPercentile"`
		DailyImprovements             int     `json:"dailyImprovements"`
		AuthorizedReplayWatched       int     `json:"authorizedReplayWatched"`
		AnonimusReplayWatched         int     `json:"anonimusReplayWatched"`
		WatchedReplays                int     `json:"watchedReplays"`
	} `json:"scoreStats"`
	LastWeekPp          float64     `json:"lastWeekPp"`
	LastWeekRank        int         `json:"lastWeekRank"`
	LastWeekCountryRank int         `json:"lastWeekCountryRank"`
	ExtensionId         int         `json:"extensionId"`
	Id                  string      `json:"id"`
	Name                string      `json:"name"`
	Platform            string      `json:"platform"`
	Avatar              string      `json:"avatar"`
	Country             string      `json:"country"`
	Alias               interface{} `json:"alias"`
	Bot                 bool        `json:"bot"`
	Pp                  float64     `json:"pp"`
	Rank                int         `json:"rank"`
	CountryRank         int         `json:"countryRank"`
	Role                string      `json:"role"`
	Socials             interface{} `json:"socials"`
	ContextExtensions   interface{} `json:"contextExtensions"`
	PatreonFeatures     interface{} `json:"patreonFeatures"`

	ClanOrder string `json:"clanOrder"`
	Clans     []struct {
		Id    int         `json:"id"`
		Tag   string      `json:"tag"`
		Color string      `json:"color"`
		Name  interface{} `json:"name"`
	} `json:"clans"`
}
type BLScoresResponse struct {
	Metadata struct {
		ItemsPerPage int `json:"itemsPerPage"`
		Page         int `json:"page"`
		Total        int `json:"total"`
	} `json:"metadata"`
	Data []struct {
		Score struct {
			Id            int     `json:"id"`
			BaseScore     int     `json:"baseScore"`
			ModifiedScore int     `json:"modifiedScore"`
			Modifiers     string  `json:"modifiers"`
			FullCombo     bool    `json:"fullCombo"`
			MaxCombo      int     `json:"maxCombo"`
			MissedNotes   int     `json:"missedNotes"`
			BadCuts       int     `json:"badCuts"`
			Hmd           int     `json:"hmd"`
			Controller    int     `json:"controller"`
			Accuracy      float64 `json:"accuracy"`
			Pp            float64 `json:"pp"`
			EpochTime     int     `json:"epochTime"`
		} `json:"score"`
		Leaderboard struct {
			Id         string `json:"id"`
			SongHash   string `json:"songHash"`
			ModeName   string `json:"modeName"`
			Difficulty int    `json:"difficulty"`
		} `json:"leaderboard"`
	} `json:"data"`
}
