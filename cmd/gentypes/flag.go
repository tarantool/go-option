package main

import (
	"maps"
	"slices"
	"strings"
)

type stringListFlag []string

func (s *stringListFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringListFlag) Set(s2 string) error {
	*s = append(*s, s2)
	return nil
}

func deleteDuplicates(s stringListFlag) stringListFlag {
	uniqMap := map[string]struct{}{}
	for _, val := range s {
		uniqMap[val] = struct{}{}
	}

	return slices.Collect(maps.Keys(uniqMap))
}

func (s *stringListFlag) Get() []string {
	deduped := deleteDuplicates(*s)
	slices.Sort(deduped)

	return deduped
}
