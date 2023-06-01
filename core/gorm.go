package core

import (
	"os"
	"telebot_v2/global"
	"telebot_v2/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	return GormMysql()
}

func RegisterTables() {
	db := global.DB

	site1 := &model.SiteVideoUrls{
		ID:            1,
		SiteName:      "东方2",
		SiteKey:       "dongfang",
		DirectPlayUrl: "https://news-zlplay.971df.com",
		CFPlayUrl:     "https://news-df2play.3d9b.com",
		CDNPlayUrl:    "https://aws-news-df2play.c20df.cc",
		VideoCover:    "https://news-dfimg.3d9b.com",
		SiteId:        2,
	}
	site2 := &model.SiteVideoUrls{
		ID:            2,
		SiteName:      "汤姆",
		SiteKey:       "tom",
		DirectPlayUrl: "https://news-tmplay.5682tom.com",
		CFPlayUrl:     "https://news-tmplay.3d9b.com",
		CDNPlayUrl:    "https://aws-tmplay3.5693tom.com",
		VideoCover:    "https://news-tmimg.3d9b.com",
		SiteId:        1,
	}
	site3 := &model.SiteVideoUrls{
		ID:            3,
		SiteName:      "好大哥",
		SiteKey:       "hdg",
		DirectPlayUrl: "https://news-hdgplay.hdg013.vip",
		CFPlayUrl:     "https://news-hdgplay.3d9b.com",
		CDNPlayUrl:    "https://aws-news-hdgplay.hdg005.cc",
		VideoCover:    "https://news-hdgpic.3d9b.com",
		SiteId:        3,
	}
	site4 := &model.SiteVideoUrls{
		ID:            4,
		SiteName:      "365导航",
		SiteKey:       "dh365",
		DirectPlayUrl: "https://news-hdgplay.hdg013.vip",
		CFPlayUrl:     "https://news-365play.3d9b.com",
		CDNPlayUrl:    "https://aws-news-hdgplay.hdg005.cc",
		VideoCover:    "https://news-365pic.3d9b.com",
		SiteId:        4,
	}
	site5 := &model.SiteVideoUrls{
		ID:            5,
		SiteName:      "柠檬",
		SiteKey:       "nm",
		DirectPlayUrl: "https://news-lemonplay.nmsp662.com",
		CFPlayUrl:     "https://news-lemonplay.3d9b.com",
		CDNPlayUrl:    "https://aws-news-lemonplay.nm016.cc",
		VideoCover:    "https://news-lemonimg.3d9b.com",
		SiteId:        23,
	}
	site6 := &model.SiteVideoUrls{
		ID:            6,
		SiteName:      "樱桃",
		SiteKey:       "yt",
		DirectPlayUrl: "https://news-cherryplay.yt1350.com",
		CFPlayUrl:     "https://news-cherryplay.3d9b.com",
		CDNPlayUrl:    "https://aws-news-cherryplay.yt026.cc",
		VideoCover:    "https://news-cherryimg.3d9b.com",
		SiteId:        22,
	}

	err := db.AutoMigrate(
		&model.SiteVideoUrls{},
	)
	if err != nil {
		global.LOG.Error("register tables failed", zap.Error(err))
		os.Exit(0)
	}

	var count int64
	db.Model(&model.SiteVideoUrls{}).Count(&count)

	if count == 0 {
		db.CreateInBatches([]model.SiteVideoUrls{
			*site1,
			*site2,
			*site3,
			*site4,
			*site5,
			*site6,
		}, 6)
	}

	global.LOG.Info("register tables success")
}
