package main

import (
	"defcon/ai"
	"defcon/config"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	go func() {

		for {

			fmt.Println("Performing checks...")

			// Refresh/Pull the config
			config := config.Generate()

			// Pull the last 100 messages from the channel using the channel ID
			messages, err := dg.ChannelMessages(config.ChannelID, 100, "", "", "")
			if err != nil {
				fmt.Println("error retrieving messages,", err)
				dg.Close()
				return
			}

			if config.CurrentTime.After(config.StartTime) && config.CurrentTime.Before(config.EndTime) {
				// If there are no prior messages, send an initial message
				if len(messages) == 0 {
					ai.FinalMessage(dg, config.ChannelID, config.DefconDaysAway)
				} else {
					// For each message that was found
					for _, message := range messages {

						// Check if the author is from the bot
						if message.Author.ID == config.ChannelID {

							// Take the message timestamp string, and map it into a time object
							lastMessageDate, err := time.Parse("2006-01-02 15:04:05.000 -0700 MST", message.Timestamp.String())
							if err != nil {
								panic(err)
							}

							// Calculate the difference in hours
							diff := config.CurrentTime.Sub(lastMessageDate)
							hours := int(diff.Hours())

							// Check if the difference is more than, or equal to 2 hours
							if hours >= 2 {
								ai.FinalMessage(dg, config.ChannelID, config.DefconDaysAway)
							}
							break
						} else {
							ai.FinalMessage(dg, config.ChannelID, config.DefconDaysAway)
							break
						}
					}
				}
			}
			time.Sleep(10 * time.Minute)
		}

	}()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
