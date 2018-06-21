package potCollector

import (
	"github.com/PuerkitoBio/goquery"
)

type Board struct {
	Threads *Threads
}

func NewBoard() *Board {
	return &Board{
		Threads: &Threads{},
	}
}

func (p *Board) LoadThreadListDocument(listURL string) (*goquery.Document, error) {
	return GetDocument(listURL)
}

func (p *Board) LoadThreads(doc *goquery.Document) error {
	var err error
	doc.Find("div > small#trad > a").Each(func(_ int, s *goquery.Selection) {
		thread := NewThread(s)
		p.Threads.Append(thread)
	})

	for _, thread := range p.Threads.List {
		if err := thread.LoadThread(); err != nil {
			return err
		}
	}

	return err
}
