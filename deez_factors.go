package main

// TODO: Add check for username or email to compare to whitelist
// TODO: Update output
// TODO: Allow for custom filter
// TODO: Add check for site owner
// TODO: Add warning if not owner of organization
/* 
> "https://api.github.com/orgs/optiv/members?filter=2fa_disabled"
{
  "message": "Only owners can use this filter.",
  "documentation_url": "https://developer.github.com/v3/orgs/members/#audit-two-factor-auth"
}
*/

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "github.com/joho/godotenv"
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
    "gopkg.in/alecthomas/kingpin.v2"
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

func checkWhiteList(name string, whitelist []string) (bool) {
    for _, value := range whitelist {
        if name == value {
            return true
        }
    }
    return false
}

func main() {
    // Command line arguments and flags
    org_name := kingpin.Flag("org", "Name of the GitHub organization.").Required().String()
    api_token := kingpin.Flag("token", "GitHub Personal API Token.").String()
    whitelist := kingpin.Flag("whitelist", "Path to whitelist file, does not print users in whitelist").String()
    use_env := kingpin.Flag("env", "Use the .env file or variable GITHUB_API_KEY.").Short('e').Bool()
    filter := kingpin.Flag("disable-filter", "Disables the user filter").Short('d').Bool()

    kingpin.Parse()

    // Use either .env file or environment variable if GITHUB_API_KEY not provided
    if *use_env {
        err := godotenv.Load()
        if err != nil {
            fmt.Println("Unable to load .env file")
        }
        *api_token = os.Getenv("GITHUB_API_KEY")
    } else if len(*api_token) == 0 {
        fmt.Println("No GITHUB_API_KEY variable found")
    }

    // If supplied, read a list of users who are allowed to have 2FA turned off
    // These users will not display on the program output
    var wlist []string
    if len(*whitelist) > 0 {
        wlist, err := readWhitelist(*whitelist)
        if wlist == nil || err != nil{
            fmt.Println("Error reading whitelist: ", err, "-- proceeding with empty whitelist")
        }
    }

    // If specified, remove the user filter, otherwise default to "2fa_disabled"
    user_filter := ""
    if !*filter {
        user_filter = "2fa_disabled"
    }

    // Authenticate to GitHub
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *api_token})
    tc := oauth2.NewClient(oauth2.NoContext, ts)

    // Create a GitHub client using the token from above
    client := github.NewClient(tc)

    // Get a list of org members that don't have 2FA enabled
    // Need to use a loop because there may be multiple pages
    // of users.
    var allUsers []*github.User
    options := &github.ListMembersOptions{Filter: user_filter}
    for {
        users, response, _ := client.Organizations.ListMembers(*org_name, options)
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
        if (len(*whitelist) > 0) {
            if checkWhiteList(*v.Login, wlist) {
            continue
        }
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