package info

import (
	"github.com/lovelaced/glitzz/config"
	"github.com/thoj/go-ircevent"
	"strings"
	"testing"
)

func TestGit(t *testing.T) {
	p := New(nil, config.Default())
	output, err := p.RunCommand(".git")
	if err != nil {
		t.Errorf("error was not nil %s", err)
	}
	if len(output) != 1 {
		t.Errorf("invalid output length %d", len(output))
	}
	if !strings.HasPrefix(output[0], "http") {
		t.Errorf("invalid output %s", output[0])
	}
}

type sender struct {
	Replies []string
}

func (s *sender) Reply(e *irc.Event, text string) {
	s.Replies = append(s.Replies, text)
}

func TestIbip(t *testing.T) {
	s := &sender{}
	p := New(s, config.Default())
	e := irc.Event{Arguments: []string{".bots"}, Code: "PRIVMSG"}
	p.HandleEvent(&e)
	if len(s.Replies) != 1 {
		t.Errorf("invalid output length %d", len(s.Replies))
	}
	if !strings.HasPrefix(s.Replies[0], "Reporting in") {
		t.Errorf("invalid output %s", s.Replies[0])
	}
}
