# Docs

## Material

1. Research paper: [hidden_lake_anonymous_network.pdf](hidden_lake_anonymous_network.pdf)
2. Presentation: [hidden_lake_anonymous_network_view.pdf](hidden_lake_anonymous_network_view.pdf)
3. Video: [// Пишем сервис с нуля в анонимной сети «Hidden Lake» //](https://www.youtube.com/watch?v=ztdtOTZqxy0)
4. Video: [Теоретически доказуемая анонимность. DC, EI и QB сети. Настройка Hidden Lake | number571](https://www.youtube.com/watch?v=o2J6ewvBKmg)

There are also English-translated versions for research paper ([hidden_lake_anonymous_network_eng.pdf](hidden_lake_anonymous_network_eng.pdf)) and presentation ([hidden_lake_anonymous_network_view_eng.pdf](hidden_lake_anonymous_network_view_eng.pdf)). The translation update may not keep up with the original papers and presentations. These translations are machine-made, so some meanings may be distorted!

## HL ports

```
1. HLS = 957x
2. HLT = 958x
3. HLM = 959x
4. HLL = 956x
5. HLE = 955x
6. HLF = 954x
7. HLR = 953x
```

## Code style

In the course of editing the project, some code styles may be added, some edited. Therefore, the current state of the project may not fully adhere to the code style, but you need to strive for it.

### 1. Prefixes

The name of the global constants must begin with the prefix 'c' (internal) or 'C' (external).
```go
const (
    cInternalConst = 1
    CExternalConst = 2
)
```

The name of the global variables must begin with the prefix 'g' (internal) or 'G' (external). The exception is errors with the prefix 'err' or 'Err'.
```go
var (
    gInternalVariable = 1
    GExternalVariable = 2
)
```

The name of the global structs must begin with the prefix 's' (internal) or 'S' (external). Also fields in the structure must begin with the prefix 'f' or 'F'.
```go
type (
    sInternalStruct struct{
        fInternalField int 
    }
    SExternalStruct struct{
        FExternalField int
    }
)
```

The name of the global interfaces must begin with the prefix 'i' (internal) or 'I' (external). Also type functions must begin with the prefix 'i' or 'I'.
```go
type (
    iInternalInterface interface{}
    IExternalInterface interface{}
)

type (
    iInternalFunc func()
    iExternalFunc func()
)
```

The name of the function parameters must begin with the prefix 'p'. Also method's object must be equal 'p'. The exception of this code style is test files.
```go
func f(pK, pV int) {}
func (p *sObject) m() {}
```

The name of the global constants, variables, structures, fields, interfaces in the test environment must begin with prefix 't' (internal) or 'T' (external).
```go
const (
    tcInternalConst = 1
    TcExternalConst = 2
)

var (
    tgInternalVariable = 1
    TgExternalVariable = 2
)

type (
    tsInternalStruct struct{
        tfInternalField int 
    }
    TsExternalStruct struct{
        TfInternalField int 
    }
)

type (
    tiInternalInterface interface{}
    TiExternalInterface interface{}
)

type (
    tiInternalFunc func()
    TiExternalFunc func()
)
```

### 2. Function / Methods names

Functions and methods should consist of two parts, where the first is a verb, the second is a noun. Standart names: Get, Set, Add, Del and etc. Example
```go
type IClient interface {
	GetIndex() (string, error)

	GetPubKey() (asymmetric.IPubKey, error)
	SetPrivKey(asymmetric.IPrivKey) error

	GetOnlines() ([]string, error)
	DelOnline(string) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(string, asymmetric.IPubKey) error
	DelFriend(string) error

	GetConnections() ([]string, error)
	AddConnection(string) error
	DelConnection(string) error

	BroadcastRequest(asymmetric.IPubKey, request.IRequest) error
	FetchRequest(asymmetric.IPubKey, request.IRequest) ([]byte, error)
}
```

### 3. If blocks

The following is allowed.
```go
if err := f(); err != nil {
    // ...
}

err := g(
    a,
    b,
    c,
)
if err != nil {
    // ...
}
```

The following is not allowed.
```go
if v {
    // ...
} else { /* eradicate the 'else' block */
    // ...
}

err := f() /* may be in if block */
if err != nil {
    // ...
}

if err := g(
    a,
    b,
    c,
); err != nil { /* not allowed multiply line-args in if block */
    // ...
}
```

### 4. Interface declaration

When a type is bound to an interface, it must be explicitly specified like this.
```go
var (
	_ types.IRunner = &sApp{}
)
```

### 5. Calling functions/methods

External simple getter functions/methods should not be used inside the package.
```go
func (p *sObject) GetSettings() ISettings {
	return p.fSettings
}
func (p *sObject) GetValue() IValue {
    p.fMutex.Lock()
    defer p.fMutex.Unlock()

    return p.fValue
}
...
func (p *sObject) DoSomething() {
	_ = p.fSettings // correct
    _ = p.GetSettings() // incorrect

    // incorrect
    p.fMutex.Lock()
    _ = p.fValue
    p.fMutex.Unlock()

    _ = p.GetValue() // correct
}
```

### 6. Args/Returns interfaces

It is not allowed to use global structures in function arguments or when returning. Interfaces should be used instead of structures.

The following is allowed.
```go
func doObject(_ IObject) {}
func newObject() IObject {}
```

The following is not allowed.
```go
func doObject(_ *SObject) {}
func newObject() *SObject {}
```
