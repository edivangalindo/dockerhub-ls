package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type ImagesResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous string  `json:"previous"`
	Results  []Image `json:"results"`
}

type Image struct {
	Name string `json:"name"`
}

type TagsResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Tag  `json:"results"`
}

type Tag struct {
	Name string `json:"name"`
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

	for _, user := range users {
		printImageNames(user)
	}

}

// Using docker api get a list of images with tags from a dockerhub user
// and print them to stdout
func printImageNames(user string) {
	api := "https://hub.docker.com/v2/repositories/" + user + "/"

	client := &http.Client{}

	var images []Image

	page := 1
	for {
		url := api + "?pagesize=100&page=" + strconv.Itoa(page)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		var imagesResponse ImagesResponse

		err = json.NewDecoder(resp.Body).Decode(&imagesResponse)
		if err != nil {
			panic(err)
		}

		images = append(images, imagesResponse.Results...)

		if imagesResponse.Next == "" {
			break
		}

		page++
	}

	page = 1
	// Get tag for each image
	for _, image := range images {
		url := api + image.Name + "/tags" + "?pagesize=100&page=" + strconv.Itoa(page)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		var tagsResponse TagsResponse

		err = json.NewDecoder(resp.Body).Decode(&tagsResponse)
		if err != nil {
			panic(err)
		}

		if tagsResponse.Results == nil {
			continue
		}

		for _, tag := range tagsResponse.Results {
			fmt.Println(user + "/" + image.Name + ":" + tag.Name)
		}

		if tagsResponse.Next == "" {
			break
		}

		page++
	}
}
