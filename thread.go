package potCollector

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Thread struct {
	Selection         *goquery.Selection
	ID                uint64
	Title             string
	LastArticleID     uint32
	LastReadArticleID uint32
	Articles          *Articles
}
type Threads struct {
	List  []*Thread
	Count uint
}

var (
	titleRep = regexp.MustCompile(`^[0-9]+: (.+) \(([0-9]+)\)$`)
	idRep    = regexp.MustCompile(`([0-9]+)/l50`)
)

func (p *Threads) Append(t *Thread) {
	p.List = append(p.List, t)
	p.Count++
}

func (p *Threads) FilterThread(f func(*Thread) bool) {
	var threadList []*Thread
	for _, thread := range p.List {
		if f(thread) {
			threadList = append(threadList, thread)
		}
	}
	p.List = threadList
	p.Count = uint(len(threadList))
}

func NewThread(s *goquery.Selection) *Thread {
	return &Thread{
		Selection: s,
		LastReadArticleID: 1,
	}
}

func (p *Thread) LoadThread() error {
	s := p.Selection

	titleData := s.First().Nodes[0].FirstChild.Data
	title := titleRep.ReplaceAllString(titleData, "$1")
	lastArticleIDStr := titleRep.ReplaceAllString(titleData, "$2")
	lastArticleID, err := strconv.ParseUint(lastArticleIDStr, 10, 64)
	if err != nil {
		return err
	}
	href, exists := s.Attr("href")
	if !exists {
		return errors.New("not fount href attribute")
	}
	idStr := idRep.ReplaceAllString(href, "$1")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	p.ID = id
	p.Title = title
	p.LastArticleID = uint32(lastArticleID)
	p.Articles = &Articles{}

	return nil
}

func (p *Thread) LoadArticleDocument(baseURL string) (*goquery.Document, error) {
	url := fmt.Sprintf("%s%d/%d-n", baseURL, p.ID, p.LastReadArticleID)
	return GetDocument(url)
}

func (p *Thread) LoadArticles(doc *goquery.Document) error {
	var err error
	doc.Find("div.thread > div.post").Each(func(_ int, s *goquery.Selection) {
		article := NewArticle(s)
		p.Articles.Append(article)
	})

	for _, article := range p.Articles.List {
		if err := article.LoadArticle(); err != nil {
			return err
		}
	}

	return err
}
