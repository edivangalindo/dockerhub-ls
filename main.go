package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type RepositoryResponse struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []Repository `json:"results"`
}

type Repository struct {
	Name               string   `json:"name"`
	Namespace          string   `json:"namespace"`
	Repository_type    string   `json:"repository_type"`
	Status             int      `json:"status"`
	Status_description string   `json:"status_description"`
	Description        string   `json:"description"`
	Is_private         bool     `json:"is_private"`
	Star_count         int      `json:"star_count"`
	Pull_count         int      `json:"pull_count"`
	Last_updated       string   `json:"last_updated"`
	Date_registered    string   `json:"date_registered"`
	Affiliation        string   `json:"affiliation"`
	Media_types        []string `json:"media_types"`
	Content_types      []string `json:"content_types"`
	Categories         []string `json:"categories"`
}

type TagResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Tag  `json:"results"`
}

type Tag struct {
	Creator               int    `json:"creator"`
	Id                    int    `json:"id"`
	Images                string `json:"images"`
	Last_updated          string `json:"last_updated"`
	Last_updater          int    `json:"last_updater"`
	Last_updater_username string `json:"last_updater_username"`
	Name                  string `json:"name"`
	Repository            int    `json:"repository"`
	Full_size             int    `json:"full_size"`
	V2                    bool   `json:"v2"`
	Tag_status            string `json:"tag_status"`
	Tag_last_pulled       string `json:"tag_last_pulled"`
	Tag_last_pushed       string `json:"tag_last_pushed"`
	Media_type            string `json:"media_type"`
	Content_type          string `json:"content_type"`
}

var silent bool

func main() {
	// Parse the flags
	lastDays := flag.Int("ld", 0, "Filter images uploaded in the last 'n' days")
	help := flag.Bool("help", false, "Show help message")
	debug := flag.Bool("debug", false, "Show debug messages")
	flag.BoolVar(&silent, "silent", false, "Suppress debug and error messages")
	flag.Parse()

	if *help {
		fmt.Println("Usage: dockerhub-ls [options]")
		fmt.Println("Options:")
		fmt.Println("  -ld int")
		fmt.Println("        Filter images uploaded in the last 'n' days (default 0, which means all images)")
		fmt.Println("  -help")
		fmt.Println("        Show this help message")
		fmt.Println("  -debug")
		fmt.Println("        Show debug messages")
		fmt.Println("  -silent")
		fmt.Println("        Suppress debug and error messages")
		return
	}

	// Check for stdin input
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "No users detected. Hint: cat users.txt | dockerhub-ls")
		os.Exit(1)
	}

	// Read stdin
	var users []string
	for {
		var user string
		_, err := fmt.Scan(&user)
		if err != nil {
			break
		}
		users = append(users, user)
	}

	for _, user := range users {
		check_user_exist(user, *lastDays, *debug)
	}
}

func get_tags(username string, repository string, lastDays int, debug bool) {
	docker_hub_tags_url := fmt.Sprintf("https://hub.docker.com/v2/namespaces/%s/repositories/%s/tags", username, repository)

	resp, err := http.Get(docker_hub_tags_url)

	if err != nil {
		log("Error checking for tags")
		return
	}

	if resp.StatusCode == 200 {
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			log("Error reading tag response")
			return
		}

		if debug {
			log(fmt.Sprintf("Response body: %s", resp_body))
		}

		var tag_resp TagResponse

		json.Unmarshal(resp_body, &tag_resp)

		if tag_resp.Count != 0 {
			for s := range tag_resp.Results {
				lastUpdated, err := time.Parse(time.RFC3339, tag_resp.Results[s].Last_updated)
				if err != nil {
					log(fmt.Sprintf("Error parsing time: %v", err))
					continue
				}
				if lastDays == 0 || time.Since(lastUpdated) <= time.Duration(lastDays)*24*time.Hour {
					fmt.Printf("%s/%s:%s\n", username, repository, tag_resp.Results[s].Name)
				}
			}
		}
		// else {
		// 	// Print the image without a tag if none exist
		// 	fmt.Printf("%s/%s\n", username, repository)
		// }
	}
}

func get_repositories(username string, lastDays int, debug bool) {
	docker_hub_repositories_url := fmt.Sprintf("https://hub.docker.com/v2/namespaces/%s/repositories", username)

	resp, err := http.Get(docker_hub_repositories_url)

	if err != nil {
		log("Error checking for repositories")
		return
	}

	if resp.StatusCode == 200 {
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			log("Error reading repository response")
			return
		}

		if debug {
			log(fmt.Sprintf("Response body: %s", resp_body))
		}

		var repo_resp RepositoryResponse

		json.Unmarshal(resp_body, &repo_resp)

		if repo_resp.Count != 0 {
			for s := range repo_resp.Results {
				get_tags(username, repo_resp.Results[s].Name, lastDays, debug)
			}
		}
	}
}

func check_user_exist(username string, lastDays int, debug bool) {
	docker_hub_user_base_url := "https://hub.docker.com/u/"

	resp, err := http.Get(docker_hub_user_base_url + username)

	if err != nil {
		log("Error checking for user existence")
		return
	}

	if resp.StatusCode == 200 {
		get_repositories(username, lastDays, debug)
	}
}

func log(message string) {
	if !silent {
		fmt.Println(message)
	}
}
