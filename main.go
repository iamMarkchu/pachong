package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gocolly/colly"
	"strconv"
	"sync"
)

const STORE_URL = "https://www.fyvor.com/stores/"

var (
	row  = 2
	lock sync.Mutex
)

type Coupon struct {
	Title       string
	Description string
	ExpireAt    string
	AddTime     string
	Verified    bool
	Code        string
	OutUrl      string
}

type Term struct {
	H1  string
	Url string
	// Coupons []Coupon
}

func main() {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "H1")
	f.SetCellValue("Sheet1", "B1", "URL")
	c := colly.NewCollector(
		// colly.Async(true),
		colly.AllowedDomains("www.fyvor.com"),
		// colly.Debugger(&debug.LogDebugger{}),
	)
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
	detailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		//coupons := make([]Coupon, 0)
		//e.ForEach("#coupon_list .c_list .ds_list", func(i int, ee *colly.HTMLElement) {
		//	coupons = append(coupons, Coupon{
		//		Title:       ee.ChildText(".coupon_title"),
		//		Description: ee.ChildText(".cpdesc"),
		//		ExpireAt:    ee.ChildText(".ex_time"),
		//		AddTime:     ee.ChildText(".add_time"),
		//		Verified:    false,
		//		Code:        "",
		//		OutUrl:      ee.Attr("id"),
		//	})
		//})
		term := Term{
			H1:  e.ChildText("h1"),
			Url: e.Request.URL.Path,
			// Coupons: coupons,
		}
		lock.Lock()
		defer lock.Unlock()
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(row), term.H1)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(row), term.Url)
		row++
		// fmt.Println(term)
	})
	//for i:= 'A'; i <= 'Z'; i++ {
	//	c.Visit(fmt.Sprintf("%s%c/", STORE_URL, i))
	//}
	c.Visit(fmt.Sprintf("%sOther", STORE_URL))

	c.Wait()
	detailCollector.Wait()
	if err := f.SaveAs("./pachong.xlsx"); err != nil {
		fmt.Println("xlsx保存出现错误：", err.Error())
	}
}
