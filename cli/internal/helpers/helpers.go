package helpers

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func GetUrlsFromFile(sourcePath string) ([]string, error) {
	contents, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	rawUrls := strings.Split(string(contents), "\n")
	urls := []string{}

	for _, url := range rawUrls {
		trimmed := strings.TrimSpace(url)
		if trimmed != "" {
			urls = append(urls, trimmed)
		}
	}

	return urls, nil
}

func GetUrlsFromStdin() ([]string, error) {
	urls := []string{}
	stdin := bufio.NewReader(os.Stdin)

	for {
		url, err := stdin.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// nothing else to read
				break
			} else {
				return nil, err
			}
		}
		urls = append(urls, strings.TrimSpace(url))
	}

	return urls, nil
}
