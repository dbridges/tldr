package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var tldrURLs = map[string]string{
	"common": "https://raw.githubusercontent.com/tldr-pages/tldr/master/pages/common/%s.md",
	"osx":    "https://raw.githubusercontent.com/tldr-pages/tldr/master/pages/osx/%s.md",
	"linux":  "https://raw.githubusercontent.com/tldr-pages/tldr/master/pages/linux/%s.md",
}

// TLDR holds the parsed data for a TLDR entry
type TLDR struct {
	Title       string
	Description []string
	Source      string
	Examples    []TLDRExample
}

func (tldr TLDR) String() string {
	tab := "    "
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(bold("NAME") + "\n")
	b.WriteString(tab + bold(tldr.Title) + " (" + tldr.Source + ")")
	b.WriteString("\n\n")
	b.WriteString(bold("DESCRIPTION") + "\n")
	for _, line := range tldr.Description {
		b.WriteString(tab + line + "\n")
	}
	b.WriteString("\n")
	b.WriteString(bold("EXAMPLES"))
	for _, example := range tldr.Examples {
		b.WriteString("\n" + tab)
		b.WriteString(example.Description)
		b.WriteString("\n\n" + tab)
		b.WriteString(tab + example.Command)
		b.WriteString("\n")
	}
	return b.String()
}

// TLDRExample holds the parsed data for a TLDR entry example
type TLDRExample struct {
	Description string
	Command     string
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	tldr(strings.Join(os.Args[1:], "-"))
}

func usage() {
	fmt.Println("usage: tldr <command_name>")
	os.Exit(0)
}

func tldr(cmd string) {
	for src, url := range tldrURLs {
		res, err := http.Get(fmt.Sprintf(url, cmd))
		must(err)
		if res.StatusCode == http.StatusNotFound {
			continue
		}
		tldr := parseTLDR(res.Body)
		tldr.Source = src
		fmt.Println(tldr)
		return
	}
	fmt.Printf("tldr for '%s' could not be found\n", cmd)
}

func parseTLDR(body io.Reader) TLDR {
	scanner := bufio.NewScanner(body)
	tldr := TLDR{Examples: []TLDRExample{}, Description: []string{}}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		switch line[0] {
		case '#':
			tldr.Title = strings.TrimPrefix(line, "# ")
		case '>':
			tldr.Description = append(tldr.Description, strings.TrimPrefix(line, "> "))
		case '-':
			tldr.Examples = append(
				tldr.Examples,
				TLDRExample{Description: strings.TrimPrefix(line, "- ")},
			)
		case '`':
			tldr.Examples[len(tldr.Examples)-1].Command = strings.Trim(line, "`")
		}
	}
	return tldr
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bold(str string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", str)
}
