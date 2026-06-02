package meshtastic

const (
	errPrefix = "pkg/network/adapters/meshtastic = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning             = &SError{"adapter running"}
	ErrClosing             = &SError{"adapter closing"}
	ErrBroadcast           = &SError{"broadcast message"}
	ErrReceiveMessage      = &SError{"receive message"}
	ErrServiceNotRunning   = &SError{"service not running"}
	ErrCreatePythonVenv    = &SError{"create python venv"}
	ErrInstallRequirements = &SError{"install requirements"}
	ErrInvalidMessageSize  = &SError{"invalid message size"}
)
