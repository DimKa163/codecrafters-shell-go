package terminal

import (
	"io"
	"strings"
	"sync"
)

type Shell struct {
	cfg          Cfg
	t            *Terminal
	wg           sync.WaitGroup
	outchan      chan string
	buf          *RuneBuffer
	inSelectMode bool
	candidates   []*CompletableItem
}

func (s *Shell) Stdout() io.Writer {
	return s.cfg.Stdout
}

func (s *Shell) Stderr() io.Writer {
	return s.cfg.Stderr
}

func NewShell(cfh Cfg, t *Terminal) *Shell {
	s := &Shell{cfg: cfh, t: t, outchan: make(chan string), buf: NewBuffer(cfh)}
	go s.exec()
	return s
}

func (s *Shell) exec() {
	s.wg.Add(1)
	defer func() {
		s.wg.Done()
		close(s.outchan)
	}()

	for {
		r := s.t.ReadRune()
		if r == 0 {
			continue
		}

		switch r {
		case CharTab:
			if s.cfg.Completer == nil || len(s.buf.buf) == 0 {
				s.t.Bell()
				continue
			}
			if len(s.candidates) == 0 {
				s.Complete()
			} else {
				s.PrintSuggestions()
				s.candidates = nil
			}
		case CharEnter, CharCtrlJ:
			s.t.WriteString("\r\n")
			data := s.buf.Reset()
			s.outchan <- string(data)
		case MetaUp:
		case MetaDown:
		case MetaBackward:
			s.buf.MoveLeft()
		case MetaForward:
			s.buf.MoveRight()
		case CharBackspace:
			s.buf.Backspace()
		default:
			s.buf.WriteRune(r)
		}
	}
}

func (s *Shell) Complete() {
	runes := s.buf.Runes()
	candidates, _ := s.cfg.Completer.Do(runes, s.buf.idx)
	if len(candidates) == 0 {
		s.t.Bell()
		return
	}
	if len(candidates) == 1 {
		c := candidates[0]
		a := append(c.value[len(runes):], ' ')
		s.buf.WriteRunes(a)
	}
	if len(candidates) > 1 {
		prefix := commonPrefix(candidates)
		if len(prefix) > len(runes) {
			s.buf.WriteRunes(prefix[len(runes):])
			return
		}
		s.t.Bell()
		s.candidates = candidates
	}
}

func commonPrefix(candidates []*CompletableItem) []rune {
	if len(candidates) == 0 {
		return nil
	}
	prefix := append([]rune{}, candidates[0].value...)
	for _, c := range candidates[1:] {
		for len(prefix) > 0 && !HasPrefix(prefix, c.value) {
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

func (s *Shell) PrintSuggestions() {
	data := s.buf.Reset()
	suggestions := make([]string, 0, len(s.candidates))
	for _, i := range s.candidates {
		suggestions = append(suggestions, strings.TrimSpace(string(i.value)))
	}
	s.t.WriteString("\r\n")
	s.t.WriteStringLine(strings.Join(suggestions, "  "))
	s.buf.WriteRunes(data)
	//s.buf.Refresh(nil)
}

func (s *Shell) ReadLine() string {
	s.buf.Refresh(nil)
	l := <-s.outchan
	return l
}

func (s *Shell) WriteLine(data ...string) {
	s.t.WriteStringLine(strings.Join(data, ""))
	s.buf.Refresh(nil)
}
