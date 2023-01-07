package config

import (
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

type Settings struct {
	ApiUrl            string
	ConfigPath        string
	DownloadDirectory string
}

func NewSettings() *Settings {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}
	return &Settings{
		os.Getenv("API_URL"),
		GetDefaultConfigPath(),
		GetDownloadPath(),
	}
}

func GetDownloadPath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("WINDOWS_DOWNLOAD_PATH")
	}
	if runtime.GOOS == "darwin" {
		return os.Getenv("MAC_DOWNLOAD_PATH")
	} else {
		return os.Getenv("LINUX_DOWNLOAD_PATH")
	}
}

// Want to have this function for now, make user able to change it thru UI
func GetDefaultConfigPath() string {
	return os.Getenv("~/.config/manga-g")
}

// Generate default yaml config file for the app settings
// func GenerateConfig() {
// 	viper.SetConfigName("config")
// 	viper.AddConfigPath(".")
// 	viper.SetConfigType("yaml")
// 	viper.Set("api_url", "http://manga-api.bytecats.codes/")
// 	viper.Set("windows_download_path", "C:\\Users\\username\\Downloads\\")
// 	viper.Set("linux_download_path", "/home/username/Downloads/")
// 	viper.Set("mac_download_path", "/Users/username/Downloads/")

// 	viper.Set("config_path", "~/.config/manga-g/config.yaml")
// 	viper.Set("config_dir", "~/.config/manga-g/")
// 	viper.Set("config_name", "config.yaml")

// 	viper.WriteConfig()
// }
