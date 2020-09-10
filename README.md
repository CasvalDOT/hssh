# HSSH 
> An heply utility to easily search and connect into the heply servers

## Dependencies
Hssh has the following dependencies:
- [fzf](https://github.com/junegunn/fzf)

## Install

Two methods are available:

### Manual
Clone the repository and then copy the binary file into a valid executable path. (Check your PATH env)
The binary is included into:

- `binaries/linux/hssh` if you have linux
- `binaries/macos/hssh` if you have MacOS

### Brew
Due to private nature of the repository, you must create first a pernonal access token
in gitlab settings tab. Copy the generated token and then export your env in your local machine

`export GITLAB_HOMEBREW_TOKEN=<Your Token>`

After that execute the following commands

```
brew tap heply/tools git@gitlab.com:Casval/homebrew-heply.git
brew install hssh
```

---

### Notes

Hssh is written in go lang, so if you have a different system you can try to compile your own binary running
the golang build command like:
`go build <path to your file>`



## Usage
To see available options and usage run:
`hssh -h`

## Aliases

Here you have three example of hssh aliases
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


