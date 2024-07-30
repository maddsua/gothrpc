package proc

func NewProcedureStepper(segments []string) ProcedureStepper {
	return ProcedureStepper{
		segments: segments,
	}
}

type ProcedureStepper struct {
	segments []string
	cursor   int
}

func (this *ProcedureStepper) Next() (string, bool) {

	if this.cursor < 0 || this.cursor >= len(this.segments) {
		return "", false
	}

	segment := this.segments[this.cursor]
	this.cursor++

	return segment, true
}

func (this *ProcedureStepper) HasNext() bool {
	return this.cursor < len(this.segments)
}
