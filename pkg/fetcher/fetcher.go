package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	. "github-network/database/pkg/types"
)

const (
	baseUrl            = "https://api.github.com"
	ResultsPerPage int = 100
	SleepTime          = 800 * time.Millisecond
)

var (
	OAuthToken string
	client     = &http.Client{}
	logger     = log.New(os.Stdout, "[fetcher]", log.Lshortfile)
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("cannot load .env file: %s", err)
	}
	OAuthToken = os.Getenv("GITHUB_OAUTH_TOKEN")
}

func getRequest(url string) (*[]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", fmt.Sprintf("token %s", OAuthToken))

	resp, err := client.Do(req)
	if err != nil {
		logger.Printf("request failed! %s", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		logger.Println(resp.Status)
		return nil, fmt.Errorf("request failed! code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("failed in reading body! %s", err)
		return nil, err
	}
	return &body, nil
}

func FetchPullRequests(owner string, repo string, page int) (*[]PullRequest, error) {
	url := baseUrl + fmt.Sprintf(
		"/repos/%s/%s/pulls?state=all&per_page=%d&page=%d",
		owner, repo, ResultsPerPage, page,
	)

	body, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	var pullRequest *[]PullRequest
	if err := json.Unmarshal(*body, &pullRequest); err != nil {
		logger.Printf("failed in marshaling body! %s", err)
		return nil, err
	}

	return pullRequest, nil
}

func FetchIssues(owner string, repo string, page int) (*[]Issue, error) {
	url := baseUrl + fmt.Sprintf(
		"/repos/%s/%s/issues?state=all&per_page=%d&page=%d",
		owner, repo, ResultsPerPage, page,
	)

	body, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	var issues *[]Issue
	if err := json.Unmarshal(*body, &issues); err != nil {
		logger.Printf("failed in marshaling body! %s", err)
		return nil, err
	}

	return issues, nil
}

func FetchUser(username string) (*User, error) {
	url := baseUrl + fmt.Sprintf(
		"/users/%s",
		username,
	)

	body, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	var user *User
	if err := json.Unmarshal(*body, &user); err != nil {
		logger.Printf("failed in marshaling body! %s", err)
		return nil, err
	}

	return user, nil
}
