package entity

type Term struct {
	Text  string
	Freq  float64
	End   int
	Start int
	Pos   string
}

type PL struct {
	Key       string
	Value     float64
	End       int
	Start     int
	SortPara1 float64
	SortPara2 float64
}

type PLTerm struct {
	TermText string
	PL
}
