package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/dbridges/tldr/models"
	"github.com/spf13/cobra"
)

var errNotFound = fmt.Errorf("Not found")
var tldrBaseURL = "https://raw.githubusercontent.com/tldr-pages/tldr/master"
var tldrPaths = map[string]string{
	"common": "/pages/common/%s.md",
	"osx":    "/pages/osx/%s.md",
	"linux":  "/pages/linux/%s.md",
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:   "view <name>",
	Short: "View a tldr page",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := tldr(strings.Join(args, "-"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

func tldr(cmd string) error {
	tldr, err := tldrCache(cmd)
	if err != nil && err != errNotFound {
		return err
	} else if err != nil {
		tldr, err = tldrWeb(cmd)
		if err != nil {
			return err
		}
	}
	fmt.Println(formatTLDR(tldr))
	return nil
}

func tldrCache(cmd string) (*models.TLDR, error) {
	dir, err := cacheDir()
	if err != nil {
		return nil, err
	}
	for src, pth := range tldrPaths {
		searchPath := fmt.Sprintf(path.Join(dir, pth), cmd)
		f, err := os.Open(path.Join(searchPath))
		if err != nil {
			continue
		}
		tldr := parseTLDR(f)
		tldr.Source = src + ", cached"
		return &tldr, nil
	}
	return nil, errNotFound
}

func tldrWeb(cmd string) (*models.TLDR, error) {
	for src, path := range tldrPaths {
		url := tldrBaseURL + path
		res, err := http.Get(fmt.Sprintf(url, cmd))
		if err != nil {
			return nil, err
		}
		if res.StatusCode == http.StatusNotFound {
			continue
		}
		tldr := parseTLDR(res.Body)
		tldr.Source = src
		return &tldr, nil
	}
	return nil, fmt.Errorf("tldr for '%s' could not be found", cmd)
}

func parseTLDR(body io.Reader) models.TLDR {
	scanner := bufio.NewScanner(body)
	tldr := models.TLDR{Examples: []models.TLDRExample{}, Description: []string{}}
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
				models.TLDRExample{Description: strings.TrimPrefix(line, "- ")},
			)
		case '`':
			tldr.Examples[len(tldr.Examples)-1].Command = strings.Trim(line, "`")
		}
	}
	return tldr
}

func formatTLDR(tldr *models.TLDR) string {
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

func bold(str string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", str)
}
