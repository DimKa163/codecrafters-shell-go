package terminal

import (
	"bufio"
	"bytes"
	"fmt"
	"sync"
)

type RuneBuffer struct {
	length    int
	rw        sync.RWMutex
	delta     int
	cfg       *Cfg
	buf       []rune
	idx       int
	line      int
	totalLine int
}

func NewBuffer(cfg *Cfg) *RuneBuffer {
	return &RuneBuffer{cfg: cfg}
}

func (r *RuneBuffer) WriteRune(c rune) {
	r.WriteRunes([]rune{c})
}

func (r *RuneBuffer) WriteRunes(s []rune) {
	r.Refresh(func() {
		tail := append(s, r.buf[r.idx:]...)
		r.buf = append(r.buf[:r.idx], tail...)
		r.idx += len(s)
		r.length += len(s)
	})
}

func (r *RuneBuffer) MoveLeft() {
	if r.idx == 0 {
		return
	}
	r.idx--
	r.MoveLeftN(1)
	r.delta += 1
}

func (r *RuneBuffer) MoveLeftN(n int) {
	if n == 0 {
		return
	}
	buf := bufio.NewWriter(r.cfg.Stdout)
	_, _ = buf.WriteString(fmt.Sprintf("\033[%dD", n))
	_ = buf.Flush()
}

func (r *RuneBuffer) MoveRight() {
	if r.idx == len(r.buf) {
		return
	}
	r.idx++
	r.MoveRightN(1)
	r.delta -= 1
}

func (r *RuneBuffer) MoveRightN(n int) {
	if n == 0 {
		return
	}
	buf := bufio.NewWriter(r.cfg.Stdout)
	_, _ = buf.WriteString(fmt.Sprintf("\033[%dC", n))
	_ = buf.Flush()
}

func (r *RuneBuffer) Runes() []rune {
	r.rw.RLock()
	defer r.rw.RUnlock()
	newr := make([]rune, len(r.buf))
	copy(newr, r.buf)
	return newr
}

func (r *RuneBuffer) Backspace() {
	if r.idx == 0 {
		return
	}
	r.Refresh(func() {
		tail := append([]rune{}, r.buf[r.idx:]...)
		r.buf = append(r.buf[:r.idx-1], tail...)
		r.idx -= 1
	})
}

func (r *RuneBuffer) Reset() []rune {
	s := r.buf
	r.idx = 0
	r.buf = []rune{}
	r.totalLine++
	r.line++
	return s
}

func (r *RuneBuffer) Refresh(f func()) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.clean()
	if f != nil {
		f()
	}
	r.print()
}
func (r *RuneBuffer) clean() {
	buf := bufio.NewWriter(r.cfg.Stdout)
	if r.length == 0 {
		_, _ = buf.WriteString("\r\033[2K")
	} else {
		_, _ = buf.WriteString("\r\033[2K")
	}
	_ = buf.Flush()
}

func (r *RuneBuffer) print() {
	b := r.output()
	_, _ = r.cfg.Stdout.Write(b)
	r.length = len(b)
	r.MoveLeftN(r.delta)
}

func (r *RuneBuffer) output() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(string(r.cfg.Prompt))
	for _, b := range r.buf {
		buf.WriteRune(b)
	}
	return buf.Bytes()
}

func (r *RuneBuffer) promptLen() int {
	return len(r.cfg.Prompt)
}
