package ranker

type Ranker struct {
	option int
}

func NewRanker(option int) *Ranker {
	return &Ranker{option}
}

func (r *Ranker) Search(query *Page) []*ResultPage {

}
