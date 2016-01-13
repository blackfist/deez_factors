# Deez Factors
Looks for people in the GitHub organization that do not have 2-factor
authentication turned on

## How to set it up
You need to set an environment variable, `GITHUB_API_KEY`, either in your
shell or by adding it to a file called `.env` in the format `GITHUB_API_KEY:blahblah`.
The token I used was a GitHub Personal Access Token. Then you just run the program.

## Does it work?
Kind of. Maybe my personal access token might not have been the right
choice because I can't get the name and email address of the users it finds,
even though I can see that stuff through the web page. Perhaps I need to do
a more complicated OAuth thing.
