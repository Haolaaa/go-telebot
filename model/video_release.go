package model

type VideoRelease struct {
	ID      int  `mapstructure:"id"`
	SiteID  uint `mapstructure:"site_id"`
	Status  int  `mapstructure:"status"`
	VideoId int  `mapstructure:"video_id"`
}

type VideoReleaseStatus struct {
	CDNPlayUrlStatus    string `json:"cdn_play_url_status"`
	CFPlayUrlStatus     string `json:"cf_play_url_status"`
	DirectPlayUrlStatus string `json:"direct_play_url_status"`
	DownUrlStatus       string `json:"down_url_status"`
	CoverUrlStatus      string `json:"cover_url_status"`
}
