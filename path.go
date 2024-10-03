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
}

func (this *procPath) next() (string, bool) {

	if this.cursor < 0 || this.cursor >= len(this.segments) {
		return "", false
	}

	segment := this.segments[this.cursor]
	this.cursor++

	return segment, true
}

func (this *procPath) hasNext() bool {
	return this.cursor < len(this.segments)
}
