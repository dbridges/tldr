package models

// TLDR holds the parsed data for a TLDR entry
type TLDR struct {
	Title       string
	Description []string
	Source      string
	Examples    []TLDRExample
}

// TLDRExample holds the parsed data for a TLDR entry example
type TLDRExample struct {
	Description string
	Command     string
}
