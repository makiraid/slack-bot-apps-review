package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

type Review struct {
	Position        int    `json:"position"`
	ID              string `json:"id"`
	Title           string `json:"title"`
	Text            string `json:"text"`
	Rating          int    `json:"rating"`
	ReviewDate      string `json:"review_date"`
	ReviewedVersion string `json:"reviewed_version"`
	Author          Author `json:"author"`
}

type Author struct {
	Name     string `json:"name"`
	AuthorID string `json:"author_id"`
}

func GetData() ([]Review, error) {
	jsonData, err := os.ReadFile("./review.json")
	if err != nil {
		log.Fatal("Error reading file: ", err)
		return nil, err
	}

	var reviews []Review

	err = json.Unmarshal(jsonData, &reviews)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	return reviews, nil
}

func PostToSlack(review Review, botToken string, channelId string) {
	fmt.Println(botToken, channelId)
	rating := review.Rating
	title := review.Title
	message := review.Text
	author := review.Author.Name

	star, color := starsAndColorBasedOnNumber(rating)

	api := slack.New(botToken)
	attachment := slack.Attachment{
		Color: color,
		Text:  fmt.Sprintf(string("%s \n *%s* \n %s \n %s"), star, title, message, author),
	}

	channelId, timestamp, err := api.PostMessage(
		channelId,
		slack.MsgOptionText("My Awesome App Name has a new iOS review", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		log.Fatalf("%s\n", err)
	}

	log.Printf("Message successfully sent to Channel %s at %s\n", channelId, timestamp)
}

func starsAndColorBasedOnNumber(num int) (string, string) {
	if num < 0 || num > 5 {
		return "Invalid number", ""
	}

	var stars string
	var color string

	switch num {
	case 0, 1, 2:
		stars = strings.Repeat("★", num) + strings.Repeat("✩", 5-num)
		color = "#FF0000"
	case 3:
		stars = strings.Repeat("★", num) + strings.Repeat("✩", 5-num)
		color = "#FFFF00"
	case 4:
		stars = strings.Repeat("★", num) + strings.Repeat("✩", 5-num)
		color = "#90EE90"
	case 5:
		stars = strings.Repeat("★", num) + strings.Repeat("✩", 5-num)
		color = "#008000"
	}

	return stars, color
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	slackChannelId := os.Getenv("SLACK_CHANNEL_ID")

	reviews, err := GetData()
	if err != nil {
		log.Fatal(err)
	}

	for _, review := range reviews {
		PostToSlack(review, botToken, slackChannelId)
	}
}
