package cmd
import (
)

/*
Bar spec-out command generates portable specification
and uploads it to bard server. After successful upload
spec-out prints spec URL to STDOUT.

	$ barc spec-out
*/
type SpecOutCmd struct {
	*BaseSubCommand
}

