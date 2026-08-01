package terminal

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"sync"
	"syscall"
	"unicode"

	"golang.org/x/term"
)

const (
	CharLineStart = 1
	CharBackward  = 2
	CharInterrupt = 3
	CharDelete    = 4
	CharLineEnd   = 5
	CharForward   = 6
	CharBell      = 7
	CharCtrlH     = 8
	CharTab       = 9
	CharCtrlJ     = 10
	CharKill      = 11
	CharCtrlL     = 12
	CharEnter     = 13
	CharNext      = 14
	CharPrev      = 16
	CharBckSearch = 18
	CharFwdSearch = 19
	CharTranspose = 20
	CharCtrlU     = 21
	CharCtrlW     = 23
	CharCtrlY     = 25
	CharCtrlZ     = 26
	CharEsc       = 27
	CharO         = 79
	CharEscapeEx  = 91
	CharBackspace = 127
)
const (
	MetaBackward rune = -iota - 1
	MetaForward
	MetaUp
	MetaDown
)

type rawwriter struct {
	io.Writer
}

func (w *rawwriter) Write(p []byte) (n int, err error) {
	p = []byte(strings.Replace(string(p), "\n", "\r\n", -1))
	return w.Writer.Write(p)
}

type writer struct {
	current io.Writer
	prev    io.Writer
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.current.Write(p)
}

func (w *writer) SetWriter(wr io.Writer) {
	w.prev = w.current
	w.current = wr
}

func (w *writer) Reset() {
	if w.prev == nil {
		return
	}
	w.current = w.prev
	w.prev = nil
}

type Terminal struct {
	cfg      *Cfg
	wg       sync.WaitGroup
	outchan  chan rune
	oldState *term.State
	stdout   *writer
	stderr   *writer
}

func New(cfg *Cfg) *Terminal {
	t := &Terminal{
		cfg:     cfg,
		stdout:  &writer{current: &rawwriter{cfg.Stdout}},
		stderr:  &writer{current: &rawwriter{cfg.Stderr}},
		outchan: make(chan rune),
	}
	go t.exec()
	return t
}

func (t *Terminal) exec() {
	t.wg.Add(1)
	_ = t.MakeRaw()
	defer func() {
		t.wg.Done()
		_ = t.Restore()
		close(t.outchan)
	}()
	var (
		isEscape   bool
		isEscapeEx bool
	)
	buf := bufio.NewReader(t.cfg.Stdin)
	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			break
		}
		if isEscape {
			isEscape = false
			if CharEscapeEx == r {
				isEscapeEx = true
				continue
			}
		}

		if isEscapeEx {
			isEscapeEx = false
			key := readEscKey(r, buf)
			r = parseRuneEscape(key.typ)
		}
		switch r {
		case CharEsc:
			isEscape = true
		default:
			t.outchan <- r
		}
	}
}

func (t *Terminal) Bell() {
	buf := bufio.NewWriter(t.stdout)
	_, _ = buf.WriteString("\x07")
	_ = buf.Flush()
}

func (t *Terminal) MakeRaw() error {
	st, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return err
	}
	t.oldState = st
	return nil
}

func (t *Terminal) ReadRune() rune {
	r, ok := <-t.outchan
	if !ok {
		return 0
	}
	return r
}

func (t *Terminal) WriteString(s string) {
	buf := bufio.NewWriter(t.stdout)
	_, _ = buf.WriteString(s)
	_ = buf.Flush()
}

func (t *Terminal) WriteStringLine(s string) {
	buf := bufio.NewWriter(t.stdout)
	_, _ = buf.WriteString(s)
	_, _ = buf.WriteString("\n")
	_ = buf.Flush()
}

func (t *Terminal) WriteErrorStringLine(s string) {
	buf := bufio.NewWriter(t.stderr)
	_, _ = buf.WriteString(s)
	_, _ = buf.WriteString("\n")
	_ = buf.Flush()
}

func (t *Terminal) Stdout() io.Writer {
	return t.stdout
}

func (t *Terminal) Stderr() io.Writer {
	return t.stderr
}
func (t *Terminal) SetStdout(stdout io.Writer) {
	t.stdout.SetWriter(stdout)
}

func (t *Terminal) SetStderr(stderr io.Writer) {
	t.stderr.SetWriter(stderr)
}
func (t *Terminal) Reset() {
	t.stdout.Reset()
	t.stderr.Reset()
}

func (t *Terminal) Restore() error {
	return term.Restore(int(syscall.Stdin), t.oldState)
}

func parseRuneEscape(r rune) rune {
	switch r {
	case 'D':
		return MetaBackward
	case 'C':
		return MetaForward
	case 'A':
		return MetaUp
	case 'B':
		return MetaDown
	default:
		return r
	}
}

type escapeKeyPair struct {
	attr string
	typ  rune
}

func readEscKey(r rune, reader *bufio.Reader) *escapeKeyPair {
	p := escapeKeyPair{}
	buf := bytes.NewBuffer(nil)
	for {
		if r == ';' {
		} else if unicode.IsNumber(r) {
		} else {
			p.typ = r
			break
		}
		buf.WriteRune(r)
		r, _, _ = reader.ReadRune()
	}
	p.attr = buf.String()
	return &p
}
