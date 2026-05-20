package ui

import (
	"strconv"

	"rime-wanxiang-updater/internal/types"
)

type schemeChoice struct {
	key   string
	value string
}

func orderedSchemeChoices() []schemeChoice {
	choices := make([]schemeChoice, 0, len(types.SchemeMap))
	for i := 1; i <= len(types.SchemeMap); i++ {
		key := strconv.Itoa(i)
		if value, ok := types.SchemeMap[key]; ok {
			choices = append(choices, schemeChoice{key: key, value: value})
		}
	}

	return choices
}
