// nolint: errcheck
package app

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	gopeer_database "github.com/number571/go-peer/pkg/storage/database"

	"github.com/gdamore/tcell/v2"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/cmd/hlp/hlp-chat/internal/database"
	"github.com/number571/hidden-lake/cmd/hlp/hlp-chat/internal/network"
	"github.com/number571/hidden-lake/cmd/hlp/hlp-chat/internal/request"
	"github.com/number571/hidden-lake/cmd/hlp/hlp-chat/internal/settings"
	"github.com/rivo/tview"
)

var (
	_ types.IRunner = &sApp{}
)

const (
	cSendMessageTemplate = "[fuchsia]%X[white]\n%s\n[gray]%s[white]\n\n"
	cRecvMessageTeamplte = "[aqua]%X[white]\n%s\n[gray]%s[white]\n\n"
)

type sApp struct {
	fDBFilePath string
	fNetworkKey string

	fDB  database.IDatabase
	fApp *tview.Application

	fPrivKey ed25519.PrivateKey
	fChanKey asymmetric.IPrivKey
}

func NewApp(pNetworkKey string, pInputPath string) types.IRunner {
	return &sApp{
		fDBFilePath: filepath.Join(pInputPath, settings.CPathDB),
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
		channel = ""
		private = ""
	)

	form := tview.NewForm().
		AddPasswordField("[white]Channel", "", 32, '*', func(text string) { channel = text }).
		AddPasswordField("[white]Private", "", 32, '*', func(text string) { private = text })

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
				dkeySize = 3 * symmetric.CCipherKeySize
			)

			keyBuilder := keybuilder.NewKeyBuilder(1<<20, []byte("priv"))
			buildBytes := keyBuilder.Build(private, seedSize+dkeySize)

			p.fPrivKey = ed25519.NewKeyFromSeed(buildBytes[:seedSize])

			dbBytes := buildBytes[seedSize:]
			var (
				authKey = dbBytes[:symmetric.CCipherKeySize]
				encrKey = dbBytes[symmetric.CCipherKeySize : 2*symmetric.CCipherKeySize]
				hashKey = dbBytes[2*symmetric.CCipherKeySize:]
			)

			var err error
			p.fDB, err = database.NewDatabase(p.fDBFilePath, [3][]byte{authKey, encrKey, hashKey})
			if err != nil {
				panic(err)
			}
		}

		pages.AddAndSwitchToPage("chat", p.getChatPage(ctx), true)
		return event
	})

	form.SetTitle(" Authorization ").SetBorder(true)
	return form
}

func (p *sApp) rememberHash(hash []byte) bool {
	hashWithPrefix := bytes.Join([][]byte{[]byte("chat_"), hash}, []byte{})
	origDB := p.fDB.GetOrigin()
	_, err := origDB.Get(hashWithPrefix)
	if err == nil {
		return false
	}
	if !errors.Is(err, gopeer_database.ErrNotFound) {
		return false
	}
	if err := origDB.Set(hashWithPrefix, []byte{}); err != nil {
		return false
	}
	return true
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

	node := network.NewHiddenLakeChatNode(
		p.fNetworkKey,
		p.fDB.GetOrigin(),
		p.fChanKey,
		p.fPrivKey,
		func(pubKey ed25519.PublicKey, hash []byte, body string) {
			if ok := p.rememberHash(hash); !ok {
				return
			}

			msg := database.SMessage{
				FSendTime: time.Now(),
				FSender:   pubKey,
				FMessage:  body,
			}
			if err := p.fDB.Insert(channelPubKey, msg); err != nil {
				return
			}

			fmt.Fprintf(
				textView,
				cRecvMessageTeamplte,
				msg.FSender,
				msg.FMessage,
				msg.FSendTime.Format(time.DateTime),
			)
		},
	)

	go func() {
		if err := node.Run(ctx); err != nil {
			panic(err)
		}
	}()

	initText := fmt.Sprintf(
		"%s{\n\t[yellow]ED25519 public key[white]: %X\n\t[yellow]Message bytes limit[white]: %d\n}\n\n",
		strings.Join(p.getLoadMessages(channelPubKey, pubKey), ""),
		pubKey,
		node.GetMessageLimitSize(),
	)

	textView.SetText(initText)
	textView.SetFocusFunc(func() { p.fApp.SetFocus(inputField) })

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

		if request.HasNotGraphicCharacters(textToSend) {
			fmt.Fprintf(textView, "[red]%s[white]\n", "only graphic chars are required")
			return
		}

		err := node.SendMessage(ctx, textToSend)
		if err != nil {
			fmt.Fprintf(textView, "[red]%s[white]\n", "failed send message")
			return
		}

		msg := database.SMessage{FSender: pubKey, FMessage: textToSend, FSendTime: time.Now()}

		if err := p.fDB.Insert(channelPubKey, msg); err != nil {
			panic(err)
		}

		fmt.Fprintf(
			textView,
			cSendMessageTemplate,
			msg.FSender,
			msg.FMessage,
			msg.FSendTime.Format(time.DateTime),
		)
	})

	chat := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 6, false).
		AddItem(inputField, 1, 1, false)
	chat.SetBorder(true).SetTitle(" Hidden Lake ")

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
		tmpl := cRecvMessageTeamplte
		if pPubKey.Equal(msg.FSender) {
			tmpl = cSendMessageTemplate
		}
		initMsgs = append(initMsgs, fmt.Sprintf(
			tmpl,
			msg.FSender,
			msg.FMessage,
			msg.FSendTime.Format(time.DateTime),
		))
	}
	return initMsgs
}
