package utils

import "strings"

const (
	ImgUrlPrefix  = "/data/org"
	DownUrlPrefix = "/data/down"
)

func FormatUrl(url string) string {
	if url == "" {
		return ""
	}

	if strings.Contains(url, ImgUrlPrefix) {
		url = strings.Replace(url, ImgUrlPrefix, "", -1)
	} else if strings.Contains(url, DownUrlPrefix) {
		url = strings.Replace(url, DownUrlPrefix, "", -1)
	}

	return url
}
