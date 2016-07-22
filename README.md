# Deez Factors
An audit tool to look for people in a GitHub organization with the primary purpose of locating users that do not have Two-factor authentication turned on. Now you can tell your users to "Get Deez Factors!"

**Note:** In order to view users who have 2FA disabled, you must be an admin/owner of the organization for which you are checking. Refer to the GitHub API Documentation [audit-two-factor-auth](https://developer.github.com/v3/orgs/members/#audit-two-factor-auth).

**Note:** In order to check the membership of users, you must belong to that organization. Refer to the GitHub API Documentation [get-organization-membership](https://developer.github.com/v3/orgs/members/#get-organization-membership)

## Usage
```
$ ./deez_factors --help
usage: deez_factors --org=ORG [<flags>]

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
      --org=ORG              Name of the GitHub organization.
      --token=TOKEN          GitHub Personal API Token.
      --whitelist=WHITELIST  Path to user whitelist file
  -e, --env                  Use the .env file or variable GITHUB_API_KEY.
  -d, --disable-filter       Disables the user filter
```

 - Organization (Required): The name of the organzation on GitHub 
    - `--org=microsoft`
 - Token: The Personal API token used to query the GitHub API 
    - `--token=api_key`
 - Whitelist: The absolute or relative path to a whitelist file of users. Whitelist can be composed of either user names or email addresses (does not work if user's email is private)
   - `--whitelist=/Users/user/Desktop/whitelist.txt`
 - Env: Overrides `--token` and reads either the environment variable GITHUB_API_KEY or a .env file
 - Disable-Filter: Disables the '2fa_disabled' user filter. Can be used to validate output if you are not a organization owner.

### Environment Variable
If not using the `--token` flag, do one of the following and provide the `-e` flag:

```bash
export GITHUB_API_KEY=api_key
```
OR
```bash
echo "api_key" > .env
```

### Whitelist
A whitelist can be used to avoid printing out specific users either by specifying their username or email address (if public):

```
username1
user@somedomain.com
username2
```

## Examples
Providing organization without a GITHUB_API_KEY

```
$ ./deez_factors --org=microsoft
Invalid GITHUB_API_KEY value
```

User is not an owner of the organization

```
$ ./deez_factors --org=microsoft --token=api_token
Only owners can use the 2fa_disabled filter
See https://developer.github.com/v3/orgs/members/#audit-two-factor-auth
* - denotes organization admin
```

Disable the `2fa_disabled` filter so you can view users

```
$ ./deez_factors.go -d --org=microsoft --token=api_token
01: [ ] user1 (John Doe) - N/A
02: [ ] user2 (John Schmoe) - user2@gmail.com
03: [ ] user3 (N/A) - N/A
04: [ ] user4 (N/A) - N/A
05: [ ] user5 (Joe Dirt) - N/A
06: [*] user6 (Thomas the Tank) - user6@gmail.com
...
* - denotes organization admin
```

Specifying an invalid GITHUB_API_KEY

```
./deez_factors -d --org=microsoft --token=invalid_api_key
401 Unauthorized. Invalid token?
```
