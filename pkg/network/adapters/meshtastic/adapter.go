package meshtastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/storage/cache"

	_ "embed"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/meshtastic/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	internal_anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"
)

type sBinaryMessagePayload struct {
	FChannel uint8  `json:"channel"`
	FMessage []byte `json:"message"`
}

const (
	netMessageChanSize = 32
)

var (
	//go:embed service/script.py
	gScriptTemplate string
	//go:embed service/requirements.txt
	gRequirementsBody string
)

var (
	_ IMeshtasticAdapter = &sMeshtasticAdapter{}
)

type sMeshtasticAdapter struct {
	fSettings   ISettings
	fNetMsgChan chan layer1.IMessage

	fCache       cache.ICache
	fServiceAddr string

	fShortName string
	fLogger    logger.ILogger
}

func NewMeshtasticAdapter(
	pSettings ISettings,
	pCache cache.ICache,
) IMeshtasticAdapter {
	adapterSettings := pSettings.GetAdapterSettings()
	switch {
	case adapterSettings.GetMessageSizeBytes() > CLimitMessageSizeBytes:
		panic("message_size_bytes > 200")
	case adapterSettings.GetWorkSizeBits() != 0:
		panic("work_size_bits != 0")
	}
	return &sMeshtasticAdapter{
		fSettings:   pSettings,
		fNetMsgChan: make(chan layer1.IMessage, netMessageChanSize),
		fCache:      pCache,
		fLogger: logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
	}
}

func (p *sMeshtasticAdapter) WithLogger(pName string, pLogger logger.ILogger) IMeshtasticAdapter {
	p.fShortName = pName
	p.fLogger = pLogger
	return p
}

func (p *sMeshtasticAdapter) Run(pCtx context.Context) error {
	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	if err := p.createPythonVenv(ctx); err != nil {
		return errors.Join(ErrCreatePythonVenv, err)
	}
	if err := p.installPythonRequirements(ctx); err != nil {
		return errors.Join(ErrInstallRequirements, err)
	}

	go func() { _ = p.runSubscriber(ctx) }()
	if err := p.runPythonScript(ctx); err != nil {
		return errors.Join(ErrRunning, err)
	}

	return p.closePythonScript()
}

func (p *sMeshtasticAdapter) Produce(pCtx context.Context, pNetMsg layer1.IMessage) error {
	msgLen := p.fSettings.GetAdapterSettings().GetMessageSizeBytes() + layer1.CMessageHeadSize
	msgGotLen := len(pNetMsg.ToBytes())

	logBuilder := anon_logger.NewLogBuilder(p.fShortName)
	logBuilder.
		WithType(internal_anon_logger.CLogBaseSendNetworkMessage).
		WithHash(pNetMsg.GetHash()).
		WithProof(pNetMsg.GetProof()).
		WithSize(msgGotLen).
		WithConn("meshtastic")

	if uint64(msgGotLen) != msgLen {
		p.fLogger.PushWarn(logBuilder)
		return ErrInvalidMessageSize
	}

	hash := encoding.HexEncode(pNetMsg.GetHash())
	_ = p.fCache.Set(hash, []byte{})

	delay := time.Duration(0)
	if maxDelay := p.fSettings.GetMaxDelayTime(); maxDelay != 0 {
		v := random.NewRandom().GetUint64() % uint64(maxDelay.Milliseconds()) // nolint:gosec
		delay = time.Duration(v) * time.Millisecond                           // nolint:gosec
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	case <-timer.C:
	}

	httpClient := &http.Client{Timeout: p.fSettings.GetWriteTimeout()}
	_, err := api.Request(
		pCtx,
		httpClient,
		http.MethodPost,
		p.fServiceAddr,
		nil,
		&sBinaryMessagePayload{
			FChannel: p.fSettings.GetChannel(),
			FMessage: pNetMsg.GetBody(),
		},
	)
	if err != nil {
		p.fLogger.PushWarn(logBuilder)
		return err
	}

	p.fLogger.PushInfo(logBuilder)
	return nil
}

func (p *sMeshtasticAdapter) Consume(pCtx context.Context) (layer1.IMessage, error) {
	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	case msg := <-p.fNetMsgChan:
		return msg, nil
	}
}

func (p *sMeshtasticAdapter) runSubscriber(pCtx context.Context) error {
	msgSize := p.fSettings.GetAdapterSettings().GetMessageSizeBytes()

	ticker := time.NewTicker(p.fSettings.GetWatchPeriod())
	defer ticker.Stop()

	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-ticker.C:
			logBuilder := anon_logger.NewLogBuilder(p.fShortName)
			logBuilder.WithConn(p.fServiceAddr)

			rsp, err := api.Request(
				pCtx,
				&http.Client{Timeout: p.fSettings.GetReadTimeout()},
				http.MethodGet,
				p.fServiceAddr,
				nil,
				nil,
			)
			if err != nil {
				p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogBaseRecvNetworkMessage))
				continue
			}

			var msgs []*sBinaryMessagePayload
			if err := json.Unmarshal(rsp, &msgs); err != nil {
				p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogBaseRecvNetworkMessage))
				continue
			}

			for _, v := range msgs {
				if v.FChannel != p.fSettings.GetChannel() {
					continue
				}
				if uint64(len(v.FMessage)) != msgSize {
					continue
				}

				msg := layer1.NewMessage(
					layer1.NewConstructSettings(&layer1.SConstructSettings{
						FSettings: p.fSettings.GetAdapterSettings(),
					}),
					v.FMessage,
				)

				logBuilder.
					WithHash(msg.GetHash()).
					WithProof(msg.GetProof()).
					WithSize(len(msg.ToBytes()))

				p.fLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseRecvNetworkMessage))

				hash := encoding.HexEncode(msg.GetHash())
				if ok := p.fCache.Set(hash, []byte{}); !ok {
					continue
				}

				if ok := p.pushMessageToChan(msg); !ok {
					p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnMessageChanOverflow))
					continue
				}
			}
		}
	}
}

func (p *sMeshtasticAdapter) createPythonVenv(pCtx context.Context) error {
	venvPath := filepath.Join(p.fSettings.GetPath(), hla_settings.CPathVenv)
	if _, err := os.Stat(venvPath); !os.IsNotExist(err) {
		return nil
	}

	cmd := exec.CommandContext(pCtx, "python", "-m", "venv", venvPath) // nolint:gosec
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (p *sMeshtasticAdapter) installPythonRequirements(pCtx context.Context) error {
	venvPath := filepath.Join(p.fSettings.GetPath(), hla_settings.CPathVenv)
	requirementsPath := filepath.Join(p.fSettings.GetPath(), hla_settings.CPathTxt)

	if _, err := os.Stat(venvPath); os.IsNotExist(err) {
		return err
	}
	if _, err := os.Stat(requirementsPath); os.IsNotExist(err) {
		if err := os.WriteFile(requirementsPath, []byte(gRequirementsBody), 0600); err != nil {
			return err
		}
	}

	var pipBin string
	if runtime.GOOS == "windows" {
		pipBin = filepath.Join(venvPath, "Scripts/pip.exe")
	} else {
		pipBin = filepath.Join(venvPath, "bin/pip")
	}

	cmd := exec.CommandContext(pCtx, pipBin, "install", "-r", requirementsPath) // nolint:gosec
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (p *sMeshtasticAdapter) runPythonScript(pCtx context.Context) error {
	venvPath := filepath.Join(p.fSettings.GetPath(), hla_settings.CPathVenv)
	scriptPath := filepath.Join(p.fSettings.GetPath(), hla_settings.CPathPy)

	serviceAddr := p.fSettings.GetAddress()
	if serviceAddr == "" {
		port, err := p.getFreePort()
		if err != nil {
			return err
		}
		serviceAddr = fmt.Sprintf("127.0.0.1:%d", port)
	}

	scriptBody := strings.NewReplacer(
		"{{devPath}}", fmt.Sprintf(`"%s"`, p.fSettings.GetDevPath()),
		"{{srvAddr}}", fmt.Sprintf(`"%s"`, serviceAddr),
	).Replace(gScriptTemplate)

	if err := os.WriteFile(scriptPath, []byte(scriptBody), 0600); err != nil {
		return err
	}

	var pythonPath string
	if runtime.GOOS == "windows" {
		pythonPath = filepath.Join(venvPath, "Scripts/python.exe")
	} else {
		pythonPath = filepath.Join(venvPath, "bin/python")
	}

	p.fServiceAddr = fmt.Sprintf("http://%s/", serviceAddr)

	cmd := exec.CommandContext(pCtx, pythonPath, scriptPath) // nolint:gosec
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (p *sMeshtasticAdapter) closePythonScript() error {
	httpClient := &http.Client{Timeout: p.fSettings.GetWriteTimeout()}
	_, err := api.Request(
		context.Background(),
		httpClient,
		http.MethodDelete,
		p.fServiceAddr,
		nil,
		nil,
	)
	return err
}

func (p *sMeshtasticAdapter) pushMessageToChan(pMsg layer1.IMessage) bool {
	select {
	case p.fNetMsgChan <- pMsg:
		return true
	default:
		return false
	}
}

func (p *sMeshtasticAdapter) getFreePort() (uint16, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer func() { _ = listener.Close() }()
	addr := listener.Addr().(*net.TCPAddr)
	return uint16(addr.Port), nil // nolint:gosec
}
