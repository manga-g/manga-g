package config

import (
	"fmt"
	"runtime"

	"github.com/spf13/viper"
)

type Settings struct {
	ApiUrl            string
	ConfigPath        string
	DownloadDirectory string
}

func NewSettings() *Settings {
	return &Settings{
		`http://manga-api.bytecats.codes/`,
		`.config/manga-g/config.yml`,
		`.downloads/`,
	}
}

func (s *Settings) SetApiUrl(url string) {
	s.ApiUrl = url
}

func GetEnvVar(key string) string {
	// use viper
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print("Error reading config file, %s", err)
	}
	return viper.GetString(key)
}

// SetDownloadPath detect user system mac, linux, windows and set the default download path
func SetDownloadPath(settings *Settings) {
	if runtime.GOOS == "windows" {
		settings.DownloadDirectory = GetEnvVar("windows_download_path")
	}
	if runtime.GOOS == "linux" {
		settings.DownloadDirectory = GetEnvVar("linux_download_path")
	}
	if runtime.GOOS == "darwin" {
		settings.DownloadDirectory = GetEnvVar("mac_download_path")
	}
}

func SetConfigPath() string {
	if runtime.GOOS == "windows" {
		return GetEnvVar("windows_config_path")
	}
	if runtime.GOOS == "linux" {
		return GetEnvVar("linux_config_path")
	}
	if runtime.GOOS == "darwin" {
		return GetEnvVar("mac_config_path")
	}
	return ""
}

// Generate default yaml config file for the app settings
func GenerateConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.Set("api_url", "http://manga-api.bytecats.codes/")
	viper.Set("windows_download_path", "C:\\Users\\username\\Downloads\\")
	viper.Set("linux_download_path", "/home/username/Downloads/")
	viper.Set("mac_download_path", "/Users/username/Downloads/")

	viper.Set("config_path", "~/.config/manga-g/config.yaml")
	viper.Set("config_dir", "~/.config/manga-g/")
	viper.Set("config_name", "config.yaml")

	viper.WriteConfig()
}
