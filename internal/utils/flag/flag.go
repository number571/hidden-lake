package flag

import (
	"slices"
	"strings"
)

func GetBoolFlagValue(pArgs, pKeyAliases []string) bool {
	for _, arg := range pArgs {
		trimArg := strings.TrimLeft(arg, "-")
		if slices.Contains(pKeyAliases, trimArg) {
			return true
		}
	}
	return false
}

func GetFlagValue(pArgs, pKeyAliases []string, pDefault string) string {
	isNextValue := false
	for _, arg := range pArgs {
		if isNextValue {
			return arg
		}
		trimArg := strings.TrimLeft(arg, "-")
		splited := strings.Split(trimArg, "=")
		if !slices.Contains(pKeyAliases, splited[0]) {
			continue
		}
		if len(splited) == 1 {
			isNextValue = true
			continue
		}
		return strings.Join(splited[1:], "=")
	}
	if isNextValue {
		panic("args has key but value is not found")
	}
	return pDefault
}
