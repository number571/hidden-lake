package flag

import (
	"slices"
	"strings"
)

var (
	_ IFlag        = &sFlag{}
	_ IFlagBuilder = &sFlag{}
)

type sFlag struct {
	fAliases      []string
	fHasValue     bool
	fDescription  string
	fDefaultValue string
}

func NewFlagBuilder(pAliases ...string) IFlagBuilder {
	return &sFlag{fAliases: pAliases}
}

func (p *sFlag) GetAliases() []string {
	return p.fAliases
}

func (p *sFlag) GetHasValue() bool {
	return p.fHasValue
}

func (p *sFlag) GetDescription() string {
	return p.fDescription
}

func (p *sFlag) GetDefaultValue() string {
	return p.fDefaultValue
}

func (p *sFlag) WithDescription(pDescription string) IFlagBuilder {
	p.fDescription = pDescription
	return p
}

func (p *sFlag) WithDefaultValue(pDefaultValue string) IFlagBuilder {
	p.fDefaultValue = pDefaultValue
	p.fHasValue = true
	return p
}

func (p *sFlag) Build() IFlag {
	return p
}

func (p *sFlag) GetBoolValue(pArgs []string) bool {
	aliases := p.GetAliases()
	for _, arg := range pArgs {
		trimArg := strings.TrimLeft(arg, "-")
		if slices.Contains(aliases, trimArg) {
			return true
		}
	}
	return false
}

func (p *sFlag) GetStringValue(pArgs []string) string {
	aliases := p.GetAliases()
	isNextValue := false
	for _, arg := range pArgs {
		if isNextValue {
			return arg
		}
		trimArg := strings.TrimLeft(arg, "-")
		splited := strings.Split(trimArg, "=")
		if !slices.Contains(aliases, splited[0]) {
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
	return p.GetDefaultValue()
}
