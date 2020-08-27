# GOOGLER BOT
## Description

This bot is intended to use Instant View of Telegram, when search engines are not available for any reason.
Bot uses search engine to find websites by keywords, and translate that web page to Telegra.ph.
Parsing is done by [mercury parser](https://github.com/postlight/mercury-parser-api), which extrats meaningful content from website.

## Getting Started

### Dependencies

Golang dependencies installed by
```bash
go get -v *link_to_library*
```
* [tgbotapi](gopkg.in/telegram-bot-api.v4)
* [googlesearch](github.com/rocketlaunchr/google-search)
* [telegraph](gitlab.com/toby3d/telegraph)

### Installing

* _tokens.go_ file should be declared with following variables: _token_, _telegraphToken_, _shortName_, _authorName_
* Binary is built by command 
```bash
go build -o googler_bot
```
* At last, docker images should be built
```bash
docker-compose build
```

### Executing program

```bash
docker-compose up
```