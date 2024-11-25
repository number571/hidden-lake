package flag

type IFlags interface {
	Get(string) IFlag
	List() []IFlag
	Validate([]string) bool
}

type IFlag interface {
	GetAliases() []string
	GetHasValue() bool
	GetDescription() string
	GetDefaultValue() string
	GetBoolValue([]string) bool
	GetStringValue([]string) string
}

type IFlagBuilder interface {
	Build() IFlag
	WithDescription(string) IFlagBuilder
	WithDefaultValue(string) IFlagBuilder
}
