# HSSH 

An utility to easily search through your ssh configs and connect into remote servers.
Also you can use a gitlab repo to store your config files and sync it!

![Screenshot](https://raw.githubusercontent.com/CasvalDOT/hssh/master/screenshot.png)

View in action!

[![asciicast](https://asciinema.org/a/L4JOn8VIieGV3EI32C9aeDCeU.svg)](https://asciinema.org/a/L4JOn8VIieGV3EI32C9aeDCeU)

## Dependencies
HSSH use fuzzysearch. So a valid binary is required.
Below you can see two examples of fuzzy finders:
- [fzf](https://github.com/junegunn/fzf) - Written in **GO**
- [skim](https://github.com/lotabout/skim) - Written in **Rust** (Not yet tested)
- [scout](https://github.com/jhbabon/scout) - Written in **Rust** (Not yet tested)

## Notes
- Hssh check ssh config files inside .ssh home folder.
- SSH config files name must contains the string **config** or must be inside a folder with a name contains string config.
For example:

```
# Allowed files
my_config
config.01
works.config.d/list
external.config/ext01
```


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

## Git
Currently the files of ssh configs can be hosted only on gitlab (next must be github). If you would use this feature
please put the files inside a directory into the root project. Ex.
```
my_gitlab_project
|
|-- my_folder
    |
    |- myconfig.01
    |- myconfig.02
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
provider: github://my_private_token:/CasvalDOT/hssh@providers
```

#### Providers
Hssh support multiple providers for fetch remote configs repository.

A provider need a connections string. The connection use the following format:
```
<driver>://<private_token>:/<repo_id>@<sub_path>
```

- **<driver>** must be **github** or **gitlab**
- **<private_token>** is the access token created for access to private repo. You can leave empty if your repo is public (ex. github://:/...).
- **<repo_id>** is the project ID for gitlab (ex. 1235667883893298) or the combo `<owner>/<repo_name>` for github (ex. CasvalDOT/hssh)
- **<sub_path>** is the subfolder where to find the files

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
