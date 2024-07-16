package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/net/html"
)

type Transcript struct {
	Text     string
	Start    float64
	Duration float64
}

// GetTranscript fetches and parses the transcript for the given videoId
func GetTranscript(videoId string) (*string, error) {
	transcripts, err := fetchTranscript(videoId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transcript: %w", err)
	}

	var text strings.Builder
	for _, t := range transcripts {
		text.WriteString(t.Text)
		text.WriteString(" ")
	}

	transcriptString := text.String()
	return &transcriptString, nil

}

// fetchTranscript fetches the transcript from a remote server
func fetchTranscript(videoId string) ([]Transcript, error) {
	baseURL := "https://youtubetranscript.com"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	q := u.Query()
	q.Set("server_vid2", videoId)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transcript: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch transcript: unexpected status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return parseTranscriptHTML(body)
}

// parseTranscriptHTML parses the HTML response and extracts transcripts
func parseTranscriptHTML(body []byte) ([]Transcript, error) {
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var transcripts []Transcript
	var parseNode func(*html.Node)

	parseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "text" {
			var transcript Transcript
			for _, attr := range n.Attr {
				switch attr.Key {
				case "start":
					transcript.Start, _ = strconv.ParseFloat(attr.Val, 64)
				case "dur":
					transcript.Duration, _ = strconv.ParseFloat(attr.Val, 64)
				}
			}
			if n.FirstChild != nil {
				transcript.Text = n.FirstChild.Data
			}
			transcripts = append(transcripts, transcript)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c)
		}
	}

	parseNode(doc)

	if len(transcripts) == 0 {
		return nil, errors.New("no transcript found")
	}

	return transcripts, nil
}

// ValidateID validates the videoId
func ValidateID(videoId string) (bool, error) {
	baseURL := "https://video.google.com/timedtext"
	u, err := url.Parse(baseURL)
	if err != nil {
		return false, fmt.Errorf("failed to parse URL: %w", err)
	}
	q := u.Query()
	q.Set("type", "track")
	q.Set("v", videoId)
	q.Set("id", "0")
	q.Set("lang", "en")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return false, fmt.Errorf("failed to validate ID: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// ProcessTranscript processes the transcript using the OpenAI API
func ProcessTranscript(transcript *string) (*string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant. The user gives you a full video transcript and a task that asks you to summarize it. Your response should be a concise video summary",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: *transcript,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to process transcript: %w", err)
	}

	return &resp.Choices[0].Message.Content, nil
}
