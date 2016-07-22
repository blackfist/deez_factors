package main

// TODO: Add option for output CSV
// TODO: Create generator for GitHub users

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "path/filepath"
    "github.com/joho/godotenv"
    "github.com/google/go-github/github"
    "github.com/fatih/color"
    "golang.org/x/oauth2"
    "gopkg.in/alecthomas/kingpin.v2"
)

func readWhitelist(path string) ([]string, error) {
    var lines []string
    absPath, _ := filepath.Abs(path)
    file, err := os.Open(absPath)

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

func checkWhiteList(check string, whitelist []string) (bool) {
    for _, value := range whitelist {
        if strings.ToLower(check) == strings.ToLower(value) {
            return true
        }
    }
    return false
}

func main() {
    // Command line flags
    org_name := kingpin.Flag("org", "Name of the GitHub organization.").Required().String()
    api_token := kingpin.Flag("token", "GitHub Personal API Token.").String()
    whitelist := kingpin.Flag("whitelist", "Path to whitelist file, does not print users in whitelist").String()
    use_env := kingpin.Flag("env", "Use the .env file or variable GITHUB_API_KEY.").Short('e').Bool()
    filter := kingpin.Flag("disable-filter", "Disables the user filter").Short('d').Bool()

    kingpin.Parse()

    // Use either .env file or environment variable if GITHUB_API_KEY not provided
    // The flag -e overrides --token
    if *use_env {
        err := godotenv.Load()
        *api_token = os.Getenv("GITHUB_API_KEY")
        
        // Only warn the user if nothing is in .env or env
        if err != nil && len(*api_token) == 0 {
            fmt.Println("Unable to load env variable GITHUB_API_KEY")
            os.Exit(1)
        }
    } else if len(*api_token) == 0 {
        fmt.Println("Invalid GITHUB_API_KEY value")
        os.Exit(1)
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
    // Need to loop through (potentially) multiple pages
    var allUsers []*github.User
    options := &github.ListMembersOptions{Filter: user_filter}
    for {
        users, response, err := client.Organizations.ListMembers(*org_name, options)
        if strings.Contains(strings.ToLower(err.Error()), "only owners") {
            color.Yellow("Only owners can use the 2fa_disabled filter")
            color.Yellow("See https://developer.github.com/v3/orgs/members/#audit-two-factor-auth")
        }
        allUsers = append(allUsers, users...)
        if response.NextPage == 0 {
            break
        }
        options.ListOptions.Page = response.NextPage
    }

    // Loop over list of users and print login, name and email (where available)
    // Don't use golang's range counter because it will skip values for whitelisted users
    counter := 1
    for _, v := range allUsers {
        // Default values for username and email
        pubname := "N/A"
        pubmail := "N/A"
        isAdmin := "[ ]"

        // Query User API for more information on user
        user, _, _ := client.Users.Get(*v.Login)

        // Query Organization API for user membership
        membership, _, _ := client.Organizations.GetOrgMembership(*user.Login, *org_name)
        if *membership.Role == "admin" {
            isAdmin = "[*]"
        }
        
        // Check if user has a name that is public
        if user.Name != nil {
            pubname = *user.Name
        }

        // Check if user has a email that is public
        if user.Email != nil {
            pubmail = *user.Email
        }

        // If the user is whitelisted, then do not print their info
        if (len(*whitelist) > 0) {
            // This may be bad for long lists beacuse we recheck at every iteration
            wlist, err := readWhitelist(*whitelist)
            
            if wlist == nil || err != nil {
                fmt.Println("Error reading whitelist: ", err)
            } else if checkWhiteList(*user.Login, wlist) || checkWhiteList(pubmail, wlist) {
                continue
            }
        }

        fmt.Printf("%02d: %s %s (%s) - %s\n", counter, isAdmin, *user.Login, pubname, pubmail)
        counter++
    }
    fmt.Println("* - denotes organization admin")
}
