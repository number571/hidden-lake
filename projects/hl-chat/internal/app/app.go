package app

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/projects/hl-chat/internal/database"
	"github.com/rivo/tview"
)

const (
	cPrintNCharsPubKey   = 16
	cHiddenLakeChatHost  = "hidden-lake-chat"
	cSendMessageTemplate = "[fuchsia][%X][white]: %s [gray]%s[white]\n"
	cRecvMessageTeamplte = "[aqua][%X][white]: %s [gray]%s[white]\n"
)

type sApp struct {
	fDBPath     string
	fNetworkKey string

	fDB  database.IDatabase
	fApp *tview.Application

	fPrivKey ed25519.PrivateKey
	fChanKey asymmetric.IPrivKey
}

func NewApp(pNetworkKey string, pDBPath string) IApp {
	return &sApp{
		fDBPath:     pDBPath,
		fNetworkKey: pNetworkKey,
		fApp:        tview.NewApplication(),
	}
}

func (p *sApp) Run(ctx context.Context) error {
	pages := tview.NewPages()
	pages.AddAndSwitchToPage("auth", p.getAuthPage(ctx, pages), true)
	return p.fApp.SetRoot(pages, true).SetFocus(pages).Run()
}

func (p *sApp) getAuthPage(ctx context.Context, pages *tview.Pages) *tview.Form {
	var (
		private = ""
		channel = ""
	)

	form := tview.NewForm().
		AddPasswordField("[white]Private", "", 32, '*', func(text string) { private = text }).
		AddPasswordField("[white]Channel", "", 32, '*', func(text string) { channel = text })

	form.SetFieldBackgroundColor(tcell.ColorGray)

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyEnter {
			form.SetFieldBackgroundColor(tcell.ColorGray)
			return event
		}

		if channel == "" || private == "" {
			form.SetFieldBackgroundColor(tcell.ColorRed)
			return tcell.NewEventKey(tcell.KeyTab, '_', tcell.ModNone)
		}

		{
			keyBuilder := keybuilder.NewKeyBuilder(1<<20, []byte("chan"))
			seed := keyBuilder.Build(channel, asymmetric.CKeySeedSize)
			p.fChanKey = asymmetric.NewPrivKeyFromSeed(seed)
		}

		{
			const (
				seedSize = 32
				dkeySize = 2 * symmetric.CCipherKeySize
			)

			keyBuilder := keybuilder.NewKeyBuilder(1<<20, []byte("priv"))
			buildBytes := keyBuilder.Build(private, seedSize+dkeySize)

			p.fPrivKey = ed25519.NewKeyFromSeed(buildBytes[:seedSize])
			db, err := database.NewDatabase(p.fDBPath, buildBytes[seedSize:])
			if err != nil {
				panic(err)
			}

			p.fDB = db
		}

		pages.AddAndSwitchToPage("chat", p.getChatPage(ctx), true)
		return event
	})

	form.SetTitle(" Authorization ").SetBorder(true)
	return form
}

func (p *sApp) getChatPage(ctx context.Context) *tview.Flex {
	textToSend := ""
	inputField := tview.NewInputField().SetLabel(">>> ").SetChangedFunc(func(text string) {
		textToSend = text
	})
	inputField.SetLabelColor(tcell.ColorWhite)
	inputField.SetFieldBackgroundColor(tcell.ColorDefault)

	channelPubKey := p.fChanKey.GetPubKey()
	pubKey := p.fPrivKey.Public().(ed25519.PublicKey)

	textView := tview.NewTextView().
		ScrollToEnd().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			p.fApp.Draw()
		})

	textView.SetText(strings.Join(p.getLoadMessages(channelPubKey, pubKey), ""))
	textView.SetFocusFunc(func() { p.fApp.SetFocus(inputField) })

	node := p.getHLNode(p.fNetworkKey, textView)
	go func() {
		if err := node.Run(ctx); err != nil {
			panic(err)
		}
	}()

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp {
			row, column := textView.GetScrollOffset()
			textView.ScrollTo(row-1, column)
		}
		if event.Key() == tcell.KeyDown {
			row, column := textView.GetScrollOffset()
			textView.ScrollTo(row+1, column)
		}
		return event
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		defer func() {
			textToSend = ""
			inputField.SetText("")
		}()

		if textToSend == "" {
			fmt.Fprintf(textView, "[red]%s[white]\n", "non-zero text is required")
			return
		}

		if hasNotGraphicCharacters(textToSend) {
			fmt.Fprintf(textView, "[red]%s[white]\n", "only graphic chars are required")
			return
		}

		salt := random.NewRandom().GetBytes(16)
		body := []byte(textToSend)
		hash := hashing.NewHMACHasher(salt, body).ToBytes()
		sign := ed25519.Sign(p.fPrivKey, hash)

		err := node.SendRequest(
			ctx,
			channelPubKey,
			request.NewRequestBuilder().
				WithHost(cHiddenLakeChatHost).
				WithHead(map[string]string{
					"pubk": hex.EncodeToString(pubKey),
					"salt": hex.EncodeToString(salt),
					"sign": hex.EncodeToString(sign),
				}).
				WithBody(body).
				Build(),
		)
		if err != nil {
			fmt.Fprintf(textView, "[red]%s[white]: %s\n", "failed send message")
			return
		}

		msg := database.SMessage{FSender: pubKey, FMessage: textToSend, FSendTime: time.Now()}
		if err := p.fDB.Insert(channelPubKey, msg); err != nil {
			panic(err)
		}

		fmt.Fprintf(
			textView,
			cSendMessageTemplate,
			msg.FSender[:cPrintNCharsPubKey],
			msg.FMessage,
			msg.FSendTime.Format(time.DateTime),
		)
	})

	chat := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 6, false).
		AddItem(inputField, 1, 1, false)
	chat.SetBorder(true).SetTitle(" Hidden Lake Chat ")

	chat.SetFocusFunc(func() { p.fApp.SetFocus(inputField) })
	return chat
}

func (p *sApp) getLoadMessages(pChannelPubKey asymmetric.IPubKey, pPubKey ed25519.PublicKey) []string {
	initMsgs := make([]string, 0, 2048)
	msgs, err := p.fDB.Select(pChannelPubKey, 2048)
	if err != nil {
		panic(err)
	}
	for _, msg := range msgs {
		initMsg := ""
		if pPubKey.Equal(msg.FSender) {
			initMsg = fmt.Sprintf(
				cSendMessageTemplate,
				pPubKey[:cPrintNCharsPubKey],
				msg.FMessage,
				msg.FSendTime.Format(time.DateTime),
			)
		} else {
			initMsg = fmt.Sprintf(
				cRecvMessageTeamplte,
				msg.FSender[:cPrintNCharsPubKey],
				msg.FMessage,
				msg.FSendTime.Format(time.DateTime),
			)
		}
		initMsgs = append(initMsgs, initMsg)
	}
	return initMsgs
}

// echo PubKey{...} | sha384sum
func getPubKeyHash(pPubKey asymmetric.IPubKey) string {
	return hashing.NewHasher([]byte(pPubKey.ToString())).ToString()
}
