# HSSH 

An utility to easily search through your ssh configs and connect into remote servers.
Also you can use a gitlab repo to store your config files and sync it!

![Screenshot](https://raw.githubusercontent.com/CasvalDOT/hssh/master/screenshot.png)

## Dependencies
HSSH use fuzzysearch. So a valid binary is required.
Below you can see two examples of fuzzy finders:
- [fzf](https://github.com/junegunn/fzf) - Written in **GO**
- [skim](https://github.com/lotabout/skim) - Written in **Rust**
- [scout](https://github.com/jhbabon/scout) - Written in **Rust**


## Install
The available methods are:

### Clone
Clone or download the repository and then inside the folder run:

- `go mod init hssh`
- `go mod vendor`
- `go build hssh`

It generate a valid binary. Put the generated binary inside a valid binary path (Check your env `$PATH`)

### Release
Download one of the releses

### Brew
Due to private nature of the repository, you must create:
1. first a personal access token in github settings tab. 
2. Copy the generated token and then export your env in your local machine

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
In alternative you can generate new config template with the
following command:

`hssh -C`

Below an example of configuration generated:

```
fuzzysearch: "fzf"
default_provider: "gitlab"
provider:
  host: "https://gitlab.com/api/v4"
  private_token: ""
  project_id: ""
  path: ""
```

#### Providers
Hssh support multiple providers for fetch remote configs repository.

NOTE: Currently gitlab in only supported.
Under the provider section please fill the following attributes:
- `host` the gitlab api url
- `private_token` The private token use to auth in gitlab.
- `project_id` The ID of the repository where to fetch configuration files
- `path` The subfolder where search files

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
alias sshexe='hssh -e'
```

## Have you found a bug?

Please open a new issue on:

https://github.com/CasvalDOT/hssh/issues

## Support the project
If you want, you can also freely donate to fund the project development!
[donate](https://paypal.me/FGortani)

## License

Copyright (c) Fabrizio Gortani

[MIT License](http://en.wikipedia.org/wiki/MIT_License)
---
