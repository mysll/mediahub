package utils

type ProcessInfo struct {
	Replaced []string
	Ignored  []string
	Offset   []string
}

var wp = newWordProcess()

type WordProcess struct {
}

func newWordProcess() *WordProcess {
	return &WordProcess{}
}

func ProcessTitle(w string) (rw string, info ProcessInfo, err error) {
	return w, ProcessInfo{}, nil
}
