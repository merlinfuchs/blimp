package internal

import "github.com/merlinfuchs/blimp/internal/config"

type Page struct {
	Title  string     `koanf:"title"`
	Layout [][]string `koanf:"layout"`
}

func parsePagesFromConfig() ([]Page, error) {
	var pages []Page

	if err := config.K.Unmarshal("pages", &pages); err != nil {
		return nil, err
	}

	return pages, nil
}
