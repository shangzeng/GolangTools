package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/feeds"
	"time"
)

type FeedController struct {
	beego.Controller
}


func (c *FeedController) Index() {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "shangzeng blog",
		Link:        &feeds.Link{Href: "http://shangzeng.club"},
		Description: "discussion about security, photos",
		Author:      &feeds.Author{Name: "shangzeng", Email: "shang_zeng@foxmail.com"},
		Created:     now,
	}
    lennum := len(GetArctlsInfo())
	for i:=0;i<lennum;i++ {
		newitem := &feeds.Item{
			Title:       GetArctlsInfo()[i].ArcName,
			Link:        &feeds.Link{Href: "http://shangzeng.club/article/"+GetArctlsInfo()[i].Name+".html"},
			Description: GetArctlsInfo()[0].Information,
			Author:      &feeds.Author{Name: "shangzeng", Email: "shang_zeng@foxmail.com"},
			Created:     now,
		}
		feed.Items = append(feed.Items, newitem)
	}
	rss, _ := feed.ToRss()
	c.Ctx.WriteString(rss)
}


