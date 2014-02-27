package augustine

import (
	"fmt"
	"sort"
)

type Posting struct {
	docID key
	pos   int64
}

// SearchPostings searches for p in a sorted slice of Postings and returns the index as specified by
// sort.Search.
func SearchPostings(postings []Posting, p Posting) int {
	return sort.Search(len(postings), func(i int) bool {
		p2 := postings[i]
		switch {
		case p.docID < p2.docID:
			return true
		case p.docID == p2.docID:
			return p.pos <= p2.pos
		}
		return false
	})
}

type FTSIndex struct {
	keys     []string
	postings map[string][]Posting // sorted
}

func NewFTSIndex() *FTSIndex {
	return &FTSIndex{
		postings: make(map[string][]Posting),
	}
}

func (x *FTSIndex) Add(term string, p Posting) {
	i := sort.SearchStrings(x.keys, term)
	fmt.Printf("\033[01;34m>>>> i: %v\x1B[m\n", i)
	if i < len(x.keys) && x.keys[i] == term {
		// Found -- append to this key
		postings := x.postings[term]
		i := SearchPostings(postings, p)
		if i < len(postings) && postings[i] == p {
			// WTF?
			return
		}
		postings = append(postings, Posting{})
		copy(postings[i+1:], postings[i:])
		postings[i] = p
		x.postings[term] = postings
	} else {
		// Not found -- insert new key
		x.keys = append(x.keys, "")
		copy(x.keys[i+1:], x.keys[i:])
		x.keys[i] = term
		x.postings[term] = []Posting{p}
	}
}
