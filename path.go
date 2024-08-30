package gothrpc

import "strings"

func newStepper(path string) procStepper {
	pathTrimmed := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
	return procStepper{
		segments: strings.Split(pathTrimmed, "/"),
	}
}

type procStepper struct {
	segments []string
	cursor   int
}

func (this *procStepper) Next() (string, bool) {

	if this.cursor < 0 || this.cursor >= len(this.segments) {
		return "", false
	}

	segment := this.segments[this.cursor]
	this.cursor++

	return segment, true
}

func (this *procStepper) HasNext() bool {
	return this.cursor < len(this.segments)
}
