package entity

type Term struct {
	Text  string
	Freq  float64
	End   int
	Start int
}

type PL struct {
	Key   string
	Value float64
	End   int
	Start int
}

type PLTerm struct {
	TermText string
	PL
}
