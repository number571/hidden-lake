package flag

import (
	"strings"
)

var (
	_ IFlagsBuilder = &sFlagsBuilder{}
	_ IFlags        = &sFlags{}
)

type sFlags []IFlag
type sFlagsBuilder []IFlagBuilder

func NewFlagsBuilder(pArgs ...IFlagBuilder) IFlagsBuilder {
	v := sFlagsBuilder(pArgs)
	return &v
}

func (p *sFlagsBuilder) Build() IFlags {
	flags := make([]IFlag, 0, len(*p))
	for _, v := range *p {
		flags = append(flags, v.Build())
	}
	return NewFlags(flags...)
}

func NewFlags(pFlags ...IFlag) IFlags {
	mapAliases := make(map[string]struct{}, len(pFlags))
	for _, v := range pFlags {
		for _, a := range v.GetAliases() {
			if _, ok := mapAliases[a]; ok {
				panic("alias_name duplicated")
			}
			mapAliases[a] = struct{}{}
		}
	}
	v := sFlags(pFlags)
	return &v
}

func (p *sFlags) Get(pName string) IFlag {
	for _, v := range *p {
		for _, n := range v.GetAliases() {
			if n == pName {
				return v
			}
		}
	}
	panic("undefined alias_name")
}

func (p *sFlags) List() []IFlag {
	return *p
}

func (p *sFlags) Validate(pArgs []string) bool {
	appArgs := p.List()
	mapArgs := make(map[string]bool, 2*len(appArgs))
	for _, v := range appArgs {
		for _, n := range v.GetAliases() {
			mapArgs[n] = v.(*sFlag).fHasValue
		}
	}
	isNextValue := false
	for _, arg := range pArgs {
		if isNextValue {
			isNextValue = false
			continue
		}
		trimArg := strings.TrimLeft(arg, "-")
		splited := strings.Split(trimArg, "=")
		withValue, ok := mapArgs[splited[0]]
		if !ok {
			return false
		}
		if withValue && len(splited) == 1 {
			isNextValue = true
		}
		continue
	}
	return !isNextValue
}
