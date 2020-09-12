# HSSH 
> An heply utility for happy coders to easily search and connect into remote servers

## Dependencies
Hssh has the following dependencies:
- [fzf](https://github.com/junegunn/fzf) - Used to search in interactive mode


## Install
Two methods are available:

### Manual
Clone or download the repository and then inside the folder run:

- `go mod init hssh`
- `go mod vendor`
- `go build hssh`

It generate a valid binary. Put the generated binary inside a valid binary path (Check your env `$PATH`)

### Brew
Due to private nature of the repository, you must create first a personal access token
in github settings tab. Copy the generated token and then export your env in your local machine

`export HOMEBREW_GITHUB_API_TOKEN=<Your Token>`

After that execute the following commands

```
brew tap heply/tools git@gitlab.com:Casval/homebrew-heply.git
brew install hssh
```


## Configuration
You must set the following params in your configuration file.
The config file can be in `/etc/hssh/config.yml` and can be overwritten 
from `~/.config/hssh/config.yml`

#### Gitlab
Under the Gitlab section you must declare the following params.

- `host` the gitlab api url
- `private_token` The private token use to auth in gitlab.
- `project_id` The ID of the repository where to fetch configuration files
- `files` The path of configuration files separated by comma

## Usage
To see available options and usage run:
`hssh -h`

## Aliases

Here some examples of hssh aliases
that you can use

```
# To get the host files
alias sshget='hssh -s'

# To list the host connections
alias sshls='hssh -l -c'

# To list the host connections using fzf
alias sshfzf='hssh -f -l -c'

# To connect to host using host configuration
alias sshexe='hssh -le'
```


