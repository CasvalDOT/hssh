package templates

// Help template
var Help = `
----------------------------------------------------
HSSH
....................................................

## USAGE
hssh <options>

## AVAILABLE OPTIONS
-h show this message
-l return all ssh connections
-f enable fuzzysearch
-e search through entire list of connections and execute the ssh command
-s sync connections file from the repository
-C create a new exemple configuration file. (this will not overwrite your actual conf)

## CUSTOM ALIAS
For maintan compatibility with sshconfig tool you can create the following
alias:

alias sshls='hssh -l -c'
alias sshget='hssh -s'
`
