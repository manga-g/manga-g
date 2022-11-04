package config

// TODO: implement better settings management (json, .env, etc)
// TODO: make settings functions more like micromans settings

// func LoadSettings() {
// gotenv.Load()
// basedApiUrl := os.Getenv("BASED_API_URL")
// fmt.Println("BASED_API_URL:", basedApiUrl)
// if basedApiUrl == "" {
//    fmt.Println("BASED_API_URL is not set in env")
//    os.Exit(1)
// }

// mangaSaveDir := os.Getenv("MANGA_SAVE_DIR")
// fmt.Println("MANGA_SAVE_DIR:", mangaSaveDir)
// if mangaSaveDir == "" {
//    fmt.Println("MANGA_SAVE_DIR is not set in env")
//    currentDirectory, err := os.Getwd()
//    if err != nil {
//        fmt.Println("Error getting current directory:", err)
//        os.Exit(1)
//    }
//    fmt.Println("Using default" + currentDirectory)
//    mangaSaveDir = "."
// }

// port := "3000"
// basedApiUrl := "http://localhost:" + port + "/"
// }
