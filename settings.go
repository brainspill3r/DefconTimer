package config

import "time"

var startTime int
var endTime int
var defconDate string
var Token string

func init() {
	// Set your configuration values here
	startTime = 9
	endTime = 10
	defconDate = "2024-08-08"
	Token = "DISCORD TOKEN HERE"
}

type Configuration struct {
	BotID          string
	ChannelID      string
	StartTime      time.Time
	EndTime        time.Time
	CurrentTime    time.Time
	DefconDaysAway int
}

func Generate() Configuration {

	var config Configuration

	// Get the current time
	currentTime := time.Now()

	config.CurrentTime = currentTime
	config.StartTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), startTime, 0, 0, 0, currentTime.Location())
	config.EndTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), endTime, 0, 0, 0, currentTime.Location())
	config.ChannelID = "DISCORD CHANNEL ID"
	config.BotID = "DISCORD BOT ID"

	// Calculate how far away Defcon is within days
	finalDate, err := time.Parse("2006-01-02", defconDate)
	if err != nil {
		panic(err)
	}

	// Calculate the difference in days
	diff := finalDate.Sub(config.CurrentTime)

	config.DefconDaysAway = int(diff.Hours() / 24)

	return config
}
