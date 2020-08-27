package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/Zhalkhas/googler_bot/gomercury"
	googlesearch "github.com/rocketlaunchr/google-search"
	"gitlab.com/toby3d/telegraph"
)

const (
	apiURL = "https://html.duckduckgo.com/html/"
)

var (
	client     *gomercury.MercuryClient
	resultsLen = 5
	account    = telegraph.Account{
		AccessToken: telegraphToken,
		AuthorName:  authorName,
		ShortName:   shortName,
	}
)

// func getQuery(query string) []string {
// 	result := make([]string, 0)

// 	baseURL, err := url.Parse(apiURL)

// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	params := url.Values{}
// 	params.Add("q", query)
// 	baseURL.RawQuery = params.Encode()
// 	resp, err := http.Get(baseURL.String())

// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	doc, err := goquery.NewDocumentFromResponse(resp)

// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	doc.Find(".result__url").Each(func(i int, sel *goquery.Selection) {
// 		urlExtracted, err := sel.Html()
// 		if err != nil {
// 			log.Panicln(err)
// 		}
// 		result = append(result, strings.TrimFunc(urlExtracted, func(r rune) bool { return unicode.IsSpace(r) }))
// 	})
// 	if len(result) == 0 {
// 		return result
// 	}
// 	for i := range result {
// 		result[i] = "https://" + result[i]
// 	}
// 	if len(result) > resultsLen {
// 		return result[:resultsLen]
// 	}
// 	return result
// }
func getQuery(query string) []string {
	ctx := context.Background()
	results := make([]string, 0)
	searches, err := googlesearch.Search(ctx, query)
	if err != nil {
		log.Println(err)
		return results
	}
	for _, search := range searches {
		results = append(results, search.URL)
	}
	if len(results) > resultsLen {
		return results[:resultsLen]
	}
	return results
}
func parseArticle(url string) string {
	log.Println("Parsing articles")
	log.Printf("Parsing %s\n", url)

	article, err := client.Parse(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	return createArticleTelegraph(article)
}

func createArticleTelegraph(page *gomercury.MercuryDocument) string {
	log.Printf("Creating article for %s\n", page.URL)
	node, err := telegraph.ContentFormat(page.Content)
	if err != nil {
		log.Panicln(err)
	}
	var title string
	if page.Title == "" {
		title = "Sample title"
	} else {
		title = page.Title
	}
	article, err := account.CreatePage(telegraph.Page{
		AuthorName:  authorName,
		ImageURL:    page.LeadImageURL,
		Title:       title,
		Description: page.Excerpt,
		Content:     node,
		AuthorURL:   page.URL,
	}, false)
	if err != nil {
		log.Println(err)
		return ""
	}
	return article.URL
}

func main() {

	client = gomercury.New(&gomercury.MercuryConfig{
		ApiKey: "",
	})
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "number":
				newNum := update.Message.CommandArguments()
				if res, err := strconv.Atoi(newNum); err != nil {
					log.Panicln(err)
				} else {
					resultsLen = res
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Number of results set to %d\n", resultsLen))
					_, err := bot.Send(msg)
					if err != nil {
						log.Panicln(err)
					}
				}
				break
			}
		} else {

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			searchResults := getQuery(update.Message.Text)
			if len(searchResults) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nothing found ¯\\_(ツ)_/¯")
				_, err := bot.Send(msg)
				if err != nil {
					log.Panicln(err)
				}
			}
			articles := make(chan string)
			for _, url := range searchResults {
				go func(url string) {
					articles <- parseArticle(url)
				}(url)
			}
			for range searchResults {
				url := <-articles
				log.Printf("Received article %s\n", url)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, url)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
