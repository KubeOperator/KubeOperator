package migrations

import "sort"

type Migration struct {
	Version  int
	FileName string
}

type Migrations struct {
	Index      uintSlice
	Migrations map[int]Migration
}

func (i *Migrations) buildIndex() {
	i.Index = make(uintSlice, 0)
	for version := range i.Migrations {
		i.Index = append(i.Index, version)
	}
	sort.Sort(i.Index)
}

type uintSlice []int

func (s uintSlice) Len() int {
	return len(s)
}

func (s uintSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s uintSlice) Less(i, j int) bool {
	return s[i] < s[j]
}
