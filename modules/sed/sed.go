package sed


import (
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/thoj/go-ircevent"
	"regexp"
	"strings"
	"os/exec"
)

var historyLimit = 100
var historyQueue = make([]*irc.Event, historyLimit)

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &sed{
		Base:   core.NewBase("sed", sender, conf),
	}
	return rv, nil
}

type sed struct {
	core.Base
}

func (r *sed) HandleEvent(event *irc.Event) {
	if event.Code == "PRIVMSG" {
		var nicks []string
		var repl []string
		if len(historyQueue) >= historyLimit {
			historyQueue = historyQueue[1:]
		}
		historyQueue = append(historyQueue, event)
		for _, msg := range historyQueue {
			nicks = append(nicks, msg.Nick)
			repl = r.sedReplace(nicks, strings.Fields(event.Message()))
		}
		go r.processReplace(repl, event)
	}
}

func (r *sed) processReplace(repl []string, e *irc.Event) {
	text := strings.Join(repl, " ")
	r.Sender.Reply(e, text)
}


func (r *sed) sedReplace(nicks []string, arguments []string) []string {
	var replaced []string
	for _, argument := range arguments {
		if isSed(nicks, argument) {
			sd, _ := regexp.Compile("([a-zA-Z/a-zA-z0-9*/[a-zA-Z0-9*/[a-zA-Z])")
			if sd.MatchString(strings.Join(arguments, " ")) {
				for i := range historyQueue {
					replaced, err := exec.Command("sed", "-e", strings.Join(arguments, " "), historyQueue[len(historyQueue)-i-1].Message()).Output()
					if err != nil {
						r.Log.Debug("error running sed", "sed", replaced, "err", err)
					}
					}
				}
			} else {
				replaced = append(replaced, "Invalid regex")
		}
		}
	return replaced
}

func isSed(nicks []string, s string) bool {
	var indirect bool
	var tmp bool
	direct, _ := regexp.MatchString("([.+s/])", s)
	for _, nick := range nicks {
			tmp, _ = regexp.MatchString(nick + "([.+s/])", s)
			indirect = indirect || tmp
	}
	return direct || indirect
}


