package main

import (
	"log"
	"os"
	"net/http"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"time"
)

func main() {
	_ = godotenv.Load()

	imageUrl := os.Getenv("IMAGE_URL")
	imageAuthToken := os.Getenv("IMAGE_AUTH_TOKEN")
	canvasId := os.Getenv("SLACK_CANVAS_ID")
	slackToken := os.Getenv("SLACK_TOKEN")

	if imageUrl == "" {
		log.Fatalf("Missing IMAGE_URL environment variable")
	}

	if imageAuthToken == "" {
		log.Fatalf("Missing IMAGE_AUTH_TOKEN environment variable")
	}

	if canvasId == "" {
		log.Fatalf("Missing SLACK_CANVAS_ID environment variable")
	}

	if slackToken == "" {
		log.Fatalf("Missing SLACK_TOKEN environment variable")
	}

	// download image
	req, err := http.NewRequest("GET", imageUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", imageAuthToken))

	client := &http.Client{}
	image, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to download image: %v", err)
	}
	defer image.Body.Close()

	api := slack.New(slackToken)

	// upload image to slack
	fileSummary, err := api.UploadFileV2(slack.UploadFileV2Parameters{
		Filename: "image.png",
		Title:    "image",
		FileSize: int(image.ContentLength),
		Reader:   image.Body,
		// Channel:  channelId,
	})
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
	}

	fileInfo, _, _, err := api.GetFileInfo(fileSummary.ID, 1, 0)
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}

	// edit canvas
	err = api.EditCanvas(slack.EditCanvasParams{
		CanvasID: canvasId,
		Changes: []slack.CanvasChange{
			{
				Operation: "replace",
				DocumentContent: slack.DocumentContent{
					Type: "markdown",
					Markdown: fmt.Sprintf(`Updated at: %s
![%s](%s)`, time.Now().UTC().Format("15:04:05"), fileSummary.ID, fileInfo.Permalink),
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to edit canvas: %v", err)
	}
}
