package spider

type ResultPage struct {
	tfidfVector []float32
	page        *Page
}

type jsonDocList struct {
	Docs []jsonDoc
}

type jsonDoc struct {
	Id  int64
	Pos []int
}

type jsonWordList struct {
	Words []jsonWord
}

type jsonWord struct {
	Id  int
	Pos []int
}

func (j *jsonDocList) addDoc(doc jsonDoc) {
	if j.Docs == nil {
		j.Docs = make([]jsonDoc, 0, 10)
	}
	j.Docs = append(j.Docs, doc)
}
