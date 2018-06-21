package potCollector

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"
	Encoding  = "Shift-JIS"
)

func GetDocument(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	defer body.Close()

	reader, err := charset.NewReaderLabel(Encoding, body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
