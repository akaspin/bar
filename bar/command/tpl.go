package command

const pathTpl  = `
PATH is git-like pattern for restrict bar to search files in tree relative to
working directory. All patterns must use forward slashes ("/") in POSIX
notation. Single asterisk ("*") matches any symbols to next slash. Double
asterisks matches any symbols across slashes. To exclude PATH prepend it with
explamation ("!"). All patterns are must match full path.

	recursive/pattern/**
	recursive/**/with.file
	non/*/recursive
	!excluded/**
`