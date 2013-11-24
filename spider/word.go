// Author Daniel Schlaug
// Written at Hong Kong University of Science and Technology in 2013

package spider

type Word struct {
	WordID    int
	Word      string
	positions []int
}

func NewWord(word string) *Word {
	return &Word{
		WordID:    -1,
		Word:      word,
		positions: nil,
	}
}

func (w *Word) AddPositions(positions []int) {
	if len(w.positions) == 0 {
		w.positions = positions
	} else {
		for _, pos := range positions {
			w.positions = append(w.positions, pos)
		}
		//w.positions = append(w.positions, positions)
	}
}

func (w *Word) Positions() []int {
	return w.positions
}

func (w *Word) TF() int {
	return len(w.positions)
}
