package flag

import "strings"

var (
	_ IFlags = &sFlags{}
)

type sFlags []IFlag

func NewFlags(pArgs ...IFlag) IFlags {
	v := sFlags(pArgs)
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
			mapArgs[n] = v.GetHasValue()
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
		if withValue {
			if len(splited) == 1 {
				isNextValue = true
			}
			continue
		}
		continue
	}
	return !isNextValue
}
