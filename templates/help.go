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
-c enable colors
-C create a new exemple configuration file. (this will not overwrite your actual conf)
`
