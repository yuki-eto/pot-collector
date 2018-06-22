package main

import (
	"fmt"
	"strings"

	"github.com/yuki-eto/pot-collector"
)

func main() {
	// 「アイドルマスタ板」からスレッド名に「ミリシタ」が含まれるスレッドの投稿を拾ってきて表示する例
	const threadListURL = "https://krsw.5ch.net/idolmaster/subback.html"
	const threadBaseURL = "https://krsw.5ch.net/test/read.cgi/idolmaster/"
	const containsString = "ミリシタ"

	fmt.Println("loading thread list...")
	board := potCollector.NewBoard()
	doc, err := board.LoadThreadListDocument(threadListURL)
	if err != nil {
		panic(err)
	}
	if err := board.LoadThreads(doc); err != nil {
		panic(err)
	}
	fmt.Printf("%d thread found\n", board.Threads.Count)

	board.Threads.FilterThread(func(t *potCollector.Thread) bool {
		return strings.Contains(t.Title, containsString)
	})
	fmt.Printf("[filterd] %d thread found\n", board.Threads.Count)

	if board.Threads.Count == 0 {
		fmt.Println("do nothing...")
		return
	}

	fmt.Println("loading thread articles...")
	for _, thread := range board.Threads.List {
		// 差分だけを取り込みたいときは、LoadArticleDocument() を呼び出す前に
		// thread.LastReadArticleID に1より大きい数を入れるとそこからの差分を取得できる
		thread.LastReadArticleID = 1
		doc, err := thread.LoadArticleDocument(threadBaseURL)
		if err != nil {
			panic(err)
		}

		if err := thread.LoadArticles(doc); err != nil {
			panic(err)
		}

		fmt.Printf("%d: %s (%d)\n", thread.ID, thread.Title, thread.LastArticleID)
		for _, article := range thread.Articles.List {
			if article.IsOver1000 {
				continue
			}
			fmt.Printf("%d: %s %v ID:%s\n", article.ID, article.Name, article.Date, article.UID)
			fmt.Println(article.Text)
			if len(article.AnchorArticleIDs) > 0 {
				fmt.Printf("AnchorArticleIDs: %v", article.AnchorArticleIDs)
				fmt.Println()
			}

			fmt.Println()
		}
	}
}
