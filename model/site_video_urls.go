package model

type SiteVideoUrls struct {
	ParentName    string `json:"parentName" gorm:"comment:父级名称"`
	SiteName      string `json:"siteName" gorm:"comment:站点名称"`
	SiteKey       string `json:"siteKey" gorm:"comment:站点标识"`
	SiteID        int    `json:"siteId" gorm:"comment:站点ID"`
	DirectPlayUrl string `json:"directPlayUrl" gorm:"comment:直连地址"`
	CFPlayUrl     string `json:"cfPlayUrl" gorm:"comment:CF地址"`
	CDNPlayUrl    string `json:"cdnPlayUrl" gorm:"comment:CDN地址"`
	VideoCover    string `json:"videoCover" gorm:"comment:视频封面"`
	DownloadUrl   string `json:"downloadUrl" gorm:"comment:下载地址"`
}
