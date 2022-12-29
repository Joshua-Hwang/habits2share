package main

import "os"

// Contains all non-secret configs
type GlobalConfig struct {
	cached           bool
	webClientId      string
	mobileClientId   string
	port             string
	sessionFilePath  string
	accountsFilePath string
	habitsFilePath   string
	todoFilePath     string
}

var globalConfig GlobalConfig

// All environment variables we use are placed here
func GetGlobalConfig() *GlobalConfig {
	if !globalConfig.cached {
		webClientId := os.Getenv("GOOGLE_WEB_CLIENT_ID")
		mobileClientId := os.Getenv("GOOGLE_MOBILE_CLIENT_ID")

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		sessionFilePath := os.Getenv("SESSIONS_FILE")
		if sessionFilePath == "" {
			sessionFilePath = "sessions.csv"
		}
		accountsFilePath := os.Getenv("ACCOUNTS_FILE")
		if accountsFilePath == "" {
			accountsFilePath = "accounts.json"
		}
		habitsFilePath := os.Getenv("HABITS_FILE")
		if habitsFilePath == "" {
			habitsFilePath = "habits.json"
		}
		todoFilePath := os.Getenv("TODO_FILE")
		if todoFilePath == "" {
			todoFilePath = "todo.json"
		}

		globalConfig = GlobalConfig{
			cached:           true,
			webClientId:      webClientId,
			mobileClientId:   mobileClientId,
			port:             port,
			sessionFilePath:  sessionFilePath,
			accountsFilePath: accountsFilePath,
			habitsFilePath:   habitsFilePath,
			todoFilePath:     todoFilePath,
		}
	}

	return &globalConfig
}
