package gothrpc

import "strings"

func newProcPath(path string) procPath {
	pathTrimmed := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
	return procPath{
		segments: strings.Split(pathTrimmed, "/"),
	}
}

type procPath struct {
	segments []string
	cursor   int
	hydrated bool
}

func (this *procPath) next() (string, bool) {

	this.hydrated = true

	if this.cursor >= len(this.segments) {
		return "", false
	}

	segment := this.segments[this.cursor]
	this.cursor++

	return segment, true
}

func (this *procPath) hasNext() bool {
	this.hydrated = true
	return this.cursor < len(this.segments)
}

func (this *procPath) at() string {

	if !this.hydrated || len(this.segments) == 0 {
		return ""
	}

	if this.cursor >= len(this.segments) {
		return this.segments[len(this.segments)-1]
	}

	return this.segments[this.cursor]
}
