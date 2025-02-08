package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Rank struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type User struct {
	Username string `json:"username"`
	Url      string `json:"url"`
}

type UnresolvedMetadata struct {
	Issues      int `json:"issues"`
	Suggestions int `json:"suggestions"`
}

type CodeChallengeMetadata struct {
	Id                 string             `json:"id"`
	Name               string             `json:"name"`
	Slug               string             `json:"slug"`
	Category           string             `json:"category"`
	PublishedAt        string             `json:"publishedAt"`
	ApprovedAt         string             `json:"approvedAt"`
	Languages          []string           `json:"languages"`
	Url                string             `json:"url"`
	Rank               Rank               `json:"rank"`
	CreatedAt          string             `json:"createdAt"`
	CreatedBy          User               `json:"createdBy"`
	ApprovedBy         User               `json:"approvedBy"`
	Description        string             `json:"description"`
	TotalAttempts      int                `json:"totalAttempts"`
	TotalCompleted     int                `json:"totalCompleted"`
	TotalStars         int                `json:"totalStars"`
	VoteScore          int                `json:"voteScore"`
	Tags               []string           `json:"tags"`
	ContributorsWanted bool               `json:"contributorsWanted"`
	Unresolved         UnresolvedMetadata `json:"unresolved"`
}

type CodeChallenge struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	Slug               string   `json:"slug"`
	CompletedLanguages []string `json:"completedLanguages"`
	CompletedAt        string   `json:"completedAt"`
}

type CodeChallenges struct {
	TotalPages int             `json:"totalPages"`
	TotalItems int             `json:"totalItems"`
	Data       []CodeChallenge `json:"data"`
}

func main() {
	// Clone the repo locally
	if _, err := os.Stat("codewars_log"); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "https://github.com/BenIsenstein/codewars_log")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err, string(output))
		}
		fmt.Println(string(output))

		cmd = exec.Command("git", "remote", "set-url", "origin", "https://"+os.Getenv("GITHUB_TOKEN")+"@github.com/BenIsenstein/codewars_log.git")
		output, err = cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err, string(output))
		}
		fmt.Println(string(output))
	}

	// Enter the repo
	err := os.Chdir("codewars_log")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch recent code challenges from Codewars
	response, err := http.Get("https://www.codewars.com/api/v1/users/BenIsenstein/code-challenges/completed")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data CodeChallenges
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	for _, challenge := range data.Data {
		res, err := http.Get("https://www.codewars.com/api/v1/code-challenges/" + challenge.Id)

		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var metadata CodeChallengeMetadata
		if err := json.Unmarshal(body, &metadata); err != nil {
			log.Fatal(err)
		}

		for _, lang := range challenge.CompletedLanguages {
			// Make directories for new languages
			if _, err := os.Stat(lang); os.IsNotExist(err) {
				err := os.Mkdir(lang, 0755)

				if err != nil {
					log.Fatal(err)
				}
			}

			// Make files for new code challenges
			if _, err := os.Stat(filepath.Join(lang, challenge.Slug+".md")); os.IsNotExist(err) {
				file, err := os.Create(filepath.Join(lang, challenge.Slug+".md"))
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				t, err := time.Parse(time.RFC3339, challenge.CompletedAt)
				if err != nil {
					log.Fatal(err)
				}

				content := "# " + challenge.Name + "\n\n" + "[Train this kata](" + metadata.Url + ")" + "\n\n" + "Completed on " + t.UTC().Format("January 2, 2006 at 3:04:05 PM UTC") + "\n\n" + metadata.Description
				_, err = file.WriteString(content)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	// Add a README.md in each language directory
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the entries and print only directories
	for _, entry := range entries {
		if entry.IsDir() {
			readme, err := os.Create(filepath.Join(entry.Name(), "README.md"))
			if err != nil {
				log.Fatal(err)
			}
			defer readme.Close()

			dir, err := os.ReadDir(entry.Name())
			if err != nil {
				log.Fatal(err)
			}

			content := "# Training " + entry.Name() + " on Codewars\n\n" + "As of " + time.Now().UTC().Format("January 2, 2006") + ", I have trained " + fmt.Sprint(len(dir)) + " problems in " + entry.Name() + " on codewars.com."

			_, err = readme.WriteString(content)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	cmd := exec.Command("git", "add", ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err, string(output))
	}
	fmt.Println(string(output))

	cmd = exec.Command("git", "commit", "-m", "chore: update log")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err, string(output))
	}
	fmt.Println(string(output))

	cmd = exec.Command("git", "push", "-u", "origin", "HEAD")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err, string(output))
	}
	fmt.Println(string(output))
}
