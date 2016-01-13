package main

import (
  "fmt"
  "log"
  "os"
  "github.com/joho/godotenv"
  "github.com/google/go-github/github"
  "golang.org/x/oauth2"
)

func main() {
  // load environment variables from .env
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  //authenticate to github
  ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_KEY")})
  tc := oauth2.NewClient(oauth2.NoContext, ts)

  // create a github client using the token from above
  client := github.NewClient(tc)

  // Get a list of org members that don't have 2FA enabled
  options := &github.ListMembersOptions{Filter: "2fa_disabled"}
  users, _, err := client.Organizations.ListMembers("heroku", options)

  // Loop over the list of users and print their name
  // User structs store values as pointers so we need to use
  // the * to get the value
  for _, v := range users {
    fmt.Print(*v.Login, " - ")

    if v.Name != nil {
      fmt.Print(*v.Name)
    } else {
      fmt.Print("No Name")
    }

    fmt.Print(" - ")
    if v.Email != nil {
      fmt.Print(*v.Email)
    } else {
      fmt.Print("No Email")
    }

    fmt.Print("\n")
  }


}
