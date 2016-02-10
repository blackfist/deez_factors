package main

import (
	"bufio"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strings"
)

func readWhitelist(path string) ([]string, error) {
	var lines []string
	file, err := os.Open(path)

	// There might be a problem opening the file. If so,
	// return the error
	if err != nil {
		return lines, err
	}

	// No error, so make sure we close the file when we're done
	defer file.Close()

	// Now read it into an array
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			// skip lines that start with #
			continue
		}
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func checkWhiteList(name string, whitelist []string) bool {
	for _, value := range whitelist {
		if name == value {
			return true
		}
	}
	return false
}

func main() {
	// load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get organization to audit from cli
	target_organization := os.Args[1]

	// read the whitelist of user names that are allowed to have
	// 2FA turned off
	whitelist, err := readWhitelist("whitelist.txt")
	if err != nil {
		log.Println("Error reading whitelist: ", err, "-- proceeding with empty whitelist")
	}

	//authenticate to github
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_KEY")})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	// create a github client using the token from above
	client := github.NewClient(tc)

	// Get a list of org members that don't have 2FA enabled
	// Need to use a loop because there may be multiple pages
	// of users.
	var allUsers []github.User
	options := &github.ListMembersOptions{Filter: "2fa_disabled"}
	for {
		users, response, _ := client.Organizations.ListMembers(target_organization, options)
		allUsers = append(allUsers, users...)
		if response.NextPage == 0 {
			break
		}
		options.ListOptions.Page = response.NextPage
	}

	// Loop over the list of users and print their name
	// User structs store values as pointers so we need to use
	// the * to get the value

	// Also need to use a different counter than the one that
	// comes with range because otherwise when we skip
	// whitelisted rows we end up with gaps in the numbers
	counter := 1
	for _, v := range allUsers {
		// If the user is whitelisted, then move on
		if checkWhiteList(*v.Login, whitelist) {
			continue
		}
		// Try to get more information about the user
		user, _, _ := client.Users.Get(*v.Login)

		fmt.Printf("%02d: ", counter)
		fmt.Print(*v.Login, " - ")

		if user.Name != nil {
			fmt.Print(*user.Name)
		} else {
			fmt.Print("No Public Name")
		}

		fmt.Print(" - ")
		if user.Email != nil {
			fmt.Print(*user.Email)
		} else {
			fmt.Print("No Public Email")
		}

		fmt.Print("\n")
		counter++
	}

}
