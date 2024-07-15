package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

type Member struct {
	Login               string `json:"login"`
	Id                  int    `json:"id"`
	Node_id             string `json:"node_id"`
	Avatar_url          string `json:"avatar_url"`
	Gravatar_id         string `json:"gravatar_id"`
	Url                 string `json:"url"`
	Html_url            string `json:"html_url"`
	Followers_url       string `json:"followers_url"`
	Following_url       string `json:"following_url"`
	Gists_url           string `json:"gists_url"`
	Starred_url         string `json:"starred_url"`
	Subscriptions_url   string `json:"subscriptions_url"`
	Organizations_url   string `json:"organizations_url"`
	Repos_url           string `json:"repos_url"`
	Events_url          string `json:"events_url"`
	Received_events_url string `json:"received_events_url"`
	Type                string `json:"type"`
	Site_admin          string `json:"site_admin"`
}

func main() {

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

	for i, user := range users {
		if i != 0 && i%1000 == 0 {
			promptForIPChange()
		}
		check_user_exist(user)
	}

}

func promptForIPChange() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Por favor, troque de IP e pressione Enter para continuar...")
	reader.ReadString('\n')
}

func get_tags(username string, repository string) {
	docker_hub_tags_url := fmt.Sprintf("https://hub.docker.com/v2/namespaces/%s/repositories/%s/tags", username, repository)

	resp, err := http.Get(docker_hub_tags_url)

	if err != nil {
		fmt.Println("[-] Erro ao checar existencia de tags")
		return
	}

	if resp.StatusCode == 200 {
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("[-] Erro ao ler resposta da obtenção de tags")
			return
		}

		var tag_resp TagResponse

		json.Unmarshal(resp_body, &tag_resp)

		if tag_resp.Count != 0 {
			for s, _ := range tag_resp.Results {
				// Imprime a imagem com tag se existir
				fmt.Printf("%s/%s:%s\n", username, repository, tag_resp.Results[s].Name)
			}
		} else {
			// Imprime a imagem sem tag se não existir
			fmt.Printf("%s/%s\n", username, repository)
		}
	}
}

func get_repositories(username string) {
	docker_hub_repositories_url := fmt.Sprintf("https://hub.docker.com/v2/namespaces/%s/repositories", username)

	resp, err := http.Get(docker_hub_repositories_url)

	if err != nil {
		fmt.Println("[-] Erro ao checar existência de repositorios")
		return
	}

	if resp.StatusCode == 200 {
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("[-] Erro ao ler resposta da checagem de repositórios")
			return
		}

		var repo_resp RepositoryResponse

		json.Unmarshal(resp_body, &repo_resp)

		if repo_resp.Count != 0 {
			for s, _ := range repo_resp.Results {
				get_tags(username, repo_resp.Results[s].Name)
			}
		}
	}
}

func check_user_exist(username string) {

	docker_hub_user_base_url := "https://hub.docker.com/u/"

	resp, err := http.Get(docker_hub_user_base_url + username)

	if err != nil {
		fmt.Println("[-] Erro ao checar existência do usuário")
		return
	}

	if resp.StatusCode == 200 {
		get_repositories(username)
	}
}
