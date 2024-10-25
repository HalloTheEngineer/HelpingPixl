package config

import (
	"github.com/disgoorg/json"
	"log/slog"
	"os"
)

var Config Configuration

func Load() {
	file, _ := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()

	stat, _ := file.Stat()
	if stat.Size() == 0 {
		defaultConfig := &Configuration{
			Formatting: Formatting{
				InternalError:        "An internal error has occurred.",
				NotEnoughArguments:   "You need to provide more arguments.",
				NoPlayerFound:        "No player named **%s** found.",
				FetchingScoresFailed: "Fetching scores of **%s** failed.",
				AlreadySniping:       "You already started a sniping process.",
				PlaylistMsg: "Using **%d** ranked plays of **%s** and **%d** ranked plays of **%s**," +
					"\n**%d** snipe-able maps were found." +
					"\nThere are **%d** common maps where you are in the lead." +
					"\nThe process took **%dms**",
				BKUpdateDesc: "Hello there!" +
					"\nI have collected some [BurgerKing](https://burgerking.de) coupons, %d in total!" +
					"\nUnfortunately, they are only available in Germany..." +
					"\nHave fun!",
				BKCouponInfoDesc: "Hello there!" +
					"\nAll current **[BurgerKing®](https://www.burgerking.de/)** coupons are fetched every day at midnight." +
					"\nTo view the coupons, use the dropdown menu below!" +
					"\nBon appetit!",
			},
			BeatSaber: BeatSaber{
				SnipeFileDescription: "Snipe %s Playlist of %s",
				HoldFileDescription:  "Hold %s Playlist of %s",
			},
			BurgerKing: BurgerKing{
				UpdateCronSyntax: "0 23 * * *",
				NoExpirationDate: ":warning: No expiration date found!",
			},
			Discord: Discord{
				BKUpdateHookList: []string{},
			},
			WebServerAPI: WebServer{
				APIPort: "9091",
			},
		}
		bytes, _ := json.MarshalIndent(defaultConfig, "", "\t")
		_, _ = file.Write(bytes)
		slog.Info("(✓) Config file created. Please relaunch the app!")
		os.Exit(0)
		return
	}

	decoder := json.NewDecoder(file)

	configuration := Configuration{}

	err := decoder.Decode(&configuration)

	if err != nil {
		slog.Error("(✕) Your config file is corrupted, not-existing or has wrong keys.\nError: ", slog.Any("err", err))
		os.Exit(1)
	} else {
		slog.Info("(✓) Config Loaded")
	}
	Config = configuration
}

type (
	Configuration struct {
		WebServerAPI WebServer  `json:"webserver"`
		Discord      Discord    `json:"discord"`
		Formatting   Formatting `json:"formatting"`
		BeatSaber    BeatSaber  `json:"beatSaber"`
		BurgerKing   BurgerKing `json:"burgerKing"`
	}
	BeatSaber struct {
		SnipeImage           string `json:"snipeImage"`
		HoldImage            string `json:"holdImage"`
		SnipeSyncUrl         string `json:"snipeSyncUrl"`
		HoldSyncUrl          string `json:"holdSyncUrl"`
		SnipeFileDescription string `json:"snipeFileDescription"`
		HoldFileDescription  string `json:"holdFileDescription"`
	}
	Discord struct {
		BotToken         string   `json:"botToken"`
		BKUpdateHookList []string `json:"bkUpdateHooks"`
	}
	Formatting struct {
		InternalError        string `json:"internalError"`
		NotEnoughArguments   string `json:"notEnoughArguments"`
		NoPlayerFound        string `json:"noPlayerFound"`
		FetchingScoresFailed string `json:"fetchingScoresFailed"`
		AlreadySniping       string `json:"alreadySniping"`
		PlaylistMsg          string `json:"playlistMsg"`
		BKUpdateDesc         string `json:"bkUpdateDesc"`
		BKCouponInfoDesc     string `json:"bkCouponInfoDesc"`
	}
	BurgerKing struct {
		UpdateCronSyntax string `json:"updateCronSyntax"`
		NoExpirationDate string `json:"noExpirationDate"`
	}
	WebServer struct {
		EnableAPI bool   `json:"enableApi"`
		APIPort   string `json:"apiPort"`
		APIUrl    string `json:"apiUrl"`
	}
)
