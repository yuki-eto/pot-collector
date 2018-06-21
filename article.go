package potCollector

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Article struct {
	Selection  *goquery.Selection
	ID         uint32
	Name       string
	Date       *time.Time
	UID        string
	Text       string
	IsOver1000 bool
}
type Articles struct {
	List  []*Article
	Count uint
}

var (
	ahrefRep = regexp.MustCompile(`<a .*?>(.*)?</a>`)
	dateRep  = regexp.MustCompile(`\([月火水木金土日]\) ([0-9]{2}:[0-9]{2}:[0-9]{2}).[0-9]{2}`)
)

func (p *Articles) Append(a *Article) {
	p.List = append(p.List, a)
	p.Count++
}

func NewArticle(s *goquery.Selection) *Article {
	return &Article{
		ID:        0,
		Selection: s,
	}
}

func (p *Article) LoadArticle() error {
	s := p.Selection

	idStr, exists := s.Attr("id")
	if !exists {
		return errors.New("not found id attribute")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}
	p.ID = uint32(id)

	meta := s.Find("div.meta").First()
	nameElem := meta.Find(".name").First().Find("b").First()
	var name string
	if nameElem.Nodes[0].FirstChild.Type == html.TextNode {
		name = nameElem.Nodes[0].FirstChild.Data
	} else {
		name = nameElem.Find("a").First().Nodes[0].FirstChild.Data
	}
	name = strings.Trim(name, " ")
	p.Name = name

	date := meta.Find(".date").First().Nodes[0].FirstChild.Data
	if date == "Over 1000" {
		p.IsOver1000 = true
	} else {
		date = dateRep.ReplaceAllString(date, " $1 +09:00")
		dt, err := time.Parse("2006/01/02 15:04:05 -07:00", date)
		if err != nil {
			return err
		}
		p.Date = &dt
		p.IsOver1000 = false
	}

	uidChild := meta.Find(".uid").First().Nodes[0].FirstChild
	var uid string
	if uidChild != nil {
		uid = uidChild.Data[3:]
	}
	p.UID = uid

	h, err := s.Find("div.message > span.escaped").First().Html()
	if err != nil {
		return err
	}

	lines := strings.Split(h, "<br/>")
	var newLines []string
	for _, line := range lines {
		s := strings.Trim(line, " ")
		s = ahrefRep.ReplaceAllString(s, "$1")
		s = html.UnescapeString(s)
		newLines = append(newLines, s)
	}
	p.Text = strings.Join(newLines, "\n")

	return nil
}
