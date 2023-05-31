package model

type SiteVideoUrls struct {
	ID            int    `gorm:"column:id;type:int;primary_key;AUTO_INCREMENT;not null"`
	SiteName      string `gorm:"column:site_name;type:varchar(255);not null"`
	SiteKey       string `gorm:"column:site_key;type:varchar(255);not null"`
	DirectPlayUrl string `gorm:"column:direct_play_url;type:varchar(255);"`
	CFPlayUrl     string `gorm:"column:cf_play_url;type:varchar(255);"`
	CDNPlayUrl    string `gorm:"column:cdn_play_url;type:varchar(255);"`
	VideoCover    string `gorm:"column:video_cover;type:varchar(255);"`
	SiteId        int    `gorm:"column:site_id;type:int;"`
}
