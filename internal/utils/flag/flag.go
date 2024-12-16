package flag

import (
	"slices"
)

var (
	_ IFlagBuilder = &sFlag{}
	_ IFlag        = &sFlag{}
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

func (p *sFlag) GetDescription() string {
	return p.fDescription
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
		if slices.Contains(aliases, arg) {
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
		if !slices.Contains(aliases, arg) {
			continue
		}
		isNextValue = true
	}
	if isNextValue {
		panic("args has key but value is not found")
	}
	return p.fDefaultValue
}
