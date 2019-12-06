package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

const STORE_URL = "https://www.fyvor.com/stores/"

var (
	row  = 2
	row2 = 2
	lock sync.Mutex
	err error
)

type Term struct {
	SeoTitle string
	H1  string
	Url string
	Breadcrumb string
	Domain string
	SnsLink string
	Description string
	// Coupons []Coupon
}

type Coupon struct {
	Title       string
	Description string
	ExpireAt    string
	AddTime     string
	Verified    string
	Code        string
	OutUrl      string
	Type        string
	Views       string
	Recommend   string
	Exclusive   string
}

func main() {
	f := excelize.NewFile()
	f.NewSheet("Sheet2")
	// 设置 term表 表头
	f.SetCellValue("Sheet1", "A1", "SeoTitle")
	f.SetCellValue("Sheet1", "B1", "H1")
	f.SetCellValue("Sheet1", "C1", "Url")
	f.SetCellValue("Sheet1", "D1", "Breadcrumb")
	f.SetCellValue("Sheet1", "E1", "Domain")
	f.SetCellValue("Sheet1", "F1", "SnsLink")
	f.SetCellValue("Sheet1", "G1", "Description")
	// 设置 coupon表 表头
	f.SetCellValue("Sheet2", "A1", "Domain")
	f.SetCellValue("Sheet2", "B1", "Title")
	f.SetCellValue("Sheet2", "C1", "Description")
	f.SetCellValue("Sheet2", "D1", "ExpireAt")
	f.SetCellValue("Sheet2", "E1", "AddTime")
	f.SetCellValue("Sheet2", "F1", "Verified")
	f.SetCellValue("Sheet2", "G1", "Code")
	f.SetCellValue("Sheet2", "H1", "OutUrl")
	f.SetCellValue("Sheet2", "I1", "Type")
	f.SetCellValue("Sheet2", "J1", "Recommend")
	f.SetCellValue("Sheet2", "K1", "Exclusive")

	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains("www.fyvor.com"),
		colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"),
	)
	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	if err = c.Limit(&colly.LimitRule{
		DomainGlob:  "*fyvor.*",
		Parallelism: 5,
		Delay:      2 * time.Second,
	}); err != nil {
		fmt.Println("设置频率出现错误")
	}

	detailCollector := c.Clone()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnHTML(".cate_content li a", func(e *colly.HTMLElement) {
		detailCollector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})
	detailCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	detailCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	detailCollector.OnHTML("html", func(e *colly.HTMLElement) {
		coupons := make([]Coupon, 0)
		e.ForEach("#coupon_list .c_list .ds_list", func(i int, ee *colly.HTMLElement) {
			coupons = append(coupons, Coupon{
				Title:       ee.ChildText(".coupon_title"),
				Description: ee.ChildText(".cpdesc"),
				ExpireAt:    ee.ChildText(".ex_time"),
				AddTime:     ee.ChildText(".add_time"),
				Verified:    ee.ChildAttr("span.verify", "class"),
				Code:        strings.TrimPrefix(ee.Attr("id"), "cb_"),
				OutUrl:      ee.Attr("data-href"),
				Type:        ee.Attr("data-type"),
				Recommend:   ee.ChildText(".coupon_recom_tag"),
				Exclusive:   ee.ChildText(".coupon_exclu_tag"),
			})
		})
		var snsMap string
		e.ForEach(".social_link ul li", func(i int, ee *colly.HTMLElement) {
			var (
				name  = ee.ChildText("a")
				value = ee.ChildAttr("a", "href")
			)
			snsMap += name + ":" + value + ","
		})
		term := Term{
			SeoTitle: e.ChildText("title"),
			H1:  e.ChildText("h1"),
			Url: e.Request.URL.Path,
			Breadcrumb: strings.TrimPrefix(e.ChildText(".page_link_n"), "Home   Coupons  "),
			Domain: strings.TrimPrefix(e.ChildText(".golink"), "Visit "),
			SnsLink: snsMap,
			Description: e.ChildText(".store_de p"),
			// Coupons: coupons,
		}
		lock.Lock()
		defer lock.Unlock()
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(row), term.SeoTitle)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(row), term.H1)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(row), term.Url)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(row), term.Breadcrumb)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(row), term.Domain)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(row), term.SnsLink)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(row), term.Description)
		row++
		for _,coupon := range coupons {
			f.SetCellValue("Sheet2", "A"+strconv.Itoa(row2), term.Domain)
			f.SetCellValue("Sheet2", "B"+strconv.Itoa(row2), coupon.Title)
			f.SetCellValue("Sheet2", "C"+strconv.Itoa(row2), coupon.Description)
			f.SetCellValue("Sheet2", "D"+strconv.Itoa(row2), coupon.ExpireAt)
			f.SetCellValue("Sheet2", "E"+strconv.Itoa(row2), coupon.AddTime)
			f.SetCellValue("Sheet2", "F"+strconv.Itoa(row2), coupon.Verified)
			f.SetCellValue("Sheet2", "G"+strconv.Itoa(row2), coupon.Code)
			f.SetCellValue("Sheet2", "H"+strconv.Itoa(row2), coupon.OutUrl)
			f.SetCellValue("Sheet2", "I"+strconv.Itoa(row2), coupon.Type)
			f.SetCellValue("Sheet2", "J"+strconv.Itoa(row2), coupon.Recommend)
			f.SetCellValue("Sheet2", "K"+strconv.Itoa(row2), coupon.Exclusive)
			row2++
		}
	})
	for i:= 'A'; i <= 'Z'; i++ {
		c.Visit(fmt.Sprintf("%s%c/", STORE_URL, i))
	}
	c.Visit(fmt.Sprintf("%sOther", STORE_URL))

	// 测试
	// detailCollector.Visit("https://www.fyvor.com/coupons/10web.io/")
	c.Wait()
	detailCollector.Wait()
	if err := f.SaveAs("./pachong2.xlsx"); err != nil {
		fmt.Println("xlsx保存出现错误：", err.Error())
	}
}
