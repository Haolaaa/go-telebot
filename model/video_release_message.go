package model

type VideoReleaseMessage struct {
	PublishedSiteName string `json:"published_site_name"`
	VideoId           int    `json:"video_id"`
	Title             string `json:"title"`
	PlayUrl           string `json:"play_url"`
	DirectPlayUrl     string `json:"direct_play_url"`
	CFPlayUrl         string `json:"cf_play_url"`
	CDNPlayUrl        string `json:"cdn_play_url"`
	CreatedAt         string `json:"created_at"`
}

type VideoReleaseKafkaMessage struct {
	TaskName string `json:"task_name"`
	Total    int    `json:"total"`
	ErrCount int    `json:"err_count"`
	VideoReleaseMessage
	VideoReleaseStatus
}
