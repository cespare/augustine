package augustine

import "sort"

type Posting struct {
	docID key
	pos   int64
}

func (p Posting) LessEq(other Posting) bool {
	switch {
	case p.docID < other.docID:
		return true
	case p.docID == other.docID:
		return p.pos <= other.pos
	}
	return false
}

// SearchPostings searches for p in a sorted slice of Postings and returns the index as specified by
// sort.Search.
func SearchPostings(postings []Posting, p Posting) int {
	return sort.Search(len(postings), func(i int) bool { return p.LessEq(postings[i]) })
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

func (x *FTSIndex) Search(terms []string) []Posting {
	if len(terms) == 0 {
		return nil
	}
	allPostings := make([][]Posting, len(terms))
	for i, term := range terms {
		allPostings[i] = x.postings[term]
	}
	sort.Sort(byLength(allPostings))

	postings := make([]Posting, len(allPostings[0]))
	copy(postings, allPostings[0])
	for _, nextPostings := range allPostings[1:] {
		if len(postings) == 0 {
			// We've eliminated everything (perhaps at least one term without matches, or else the intersection is
			// empty)
			return nil
		}
		var i1, i2, insert int
		for i1 < len(postings) && i2 < len(nextPostings) {
			p1, p2 := postings[i1], nextPostings[i2]
			if p1.LessEq(p2) {
				if p1 == p2 {
					postings[insert] = p1
					insert++
					i1++
					i2++
					continue
				}
				i1++
				continue
			}
			i2++
			continue
		}
		postings = postings[:insert+1]
	}
	return postings
}

type byLength [][]Posting

func (b byLength) Len() int           { return len(b) }
func (b byLength) Less(i, j int) bool { return len(b[i]) < len(b[j]) }
func (b byLength) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
