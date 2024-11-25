package flag

type IFlagsBuilder interface {
	Build() IFlags
}

type IFlags interface {
	Get(string) IFlag
	List() []IFlag
	Validate([]string) bool
}

type IFlagBuilder interface {
	Build() IFlag
	WithDescription(string) IFlagBuilder
	WithDefaultValue(string) IFlagBuilder
}

type IFlag interface {
	GetAliases() []string
	GetHasValue() bool
	GetDescription() string
	GetDefaultValue() string
	GetBoolValue([]string) bool
	GetStringValue([]string) string
}
