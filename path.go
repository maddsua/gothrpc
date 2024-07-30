package gothrpc

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
