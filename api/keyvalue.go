package api

import (
	"sort"
)

type KeyValue struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func SortKeyValue(entries []KeyValue) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value < entries[j].Value
	})
}
