package main

import (
	"context"
	"image/color"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/pkg/api/kernel/client/proc"
	dto_f "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
	"github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
)

var (
	inputFriendNameEntry    *widget.Entry
	inputFriendPubKeyEntry  *widget.Entry
	inputConnectionEntry    *widget.Entry
	inputMessageEntry       *widget.Entry
	connectionSettingsLabel *widget.Label
	aboutBodyContainer      *fyne.Container
	scrollChatContainer     *container.Scroll
)

var (
	chatListenerActive = false
	closeListenChat    = make(chan struct{})
	friendNameInChat   string
)

type sConnection struct {
	online  bool
	address string
}

var (
	gFriends     = []string{}
	gConnections = []sConnection{}
)

func setAboutContent(w fyne.Window) {
	clearAfterSwitch()

	getServicesState(w)

	w.SetContent(aboutContainer)
}

func setConnectionsContent(w fyne.Window) {
	clearAfterSwitch()

	getConnections(w)
	getConnectionSettings(w)

	w.SetContent(connectionsContainer)
	w.Canvas().Focus(inputConnectionEntry)
}

func setFriendsContent(w fyne.Window) {
	clearAfterSwitch()
	getFriends(w)

	w.SetContent(friendsContainer)
	w.Canvas().Focus(inputFriendNameEntry)
}

func setChatListContent(w fyne.Window) {
	clearAfterSwitch()
	getFriends(w)

	w.SetContent(chatListContainer)
}

func setChatFriendContent(w fyne.Window, friend string) {
	clearAfterSwitch()

	friendNameInChat = friend
	loadMessages(w, friend)
	listenMessages(w, friend)

	scrollChatContainer.ScrollToBottom()
	w.SetContent(chatFriendContainer)
	w.Canvas().Focus(inputMessageEntry)
}

func clearAfterSwitch() {
	inputConnectionEntry.SetText("")
	inputFriendNameEntry.SetText("")
	inputFriendPubKeyEntry.SetText("")
	inputMessageEntry.SetText("")
	scrollChatContainer.Content.(*fyne.Container).RemoveAll()
	if chatListenerActive {
		closeListenChat <- struct{}{}
	}
}

func setServiceState(label *widget.Label, err error) {
	if err != nil {
		label.SetText("<dead>")
		label.Importance = widget.DangerImportance
		return
	}
	label.SetText("<alive>")
	label.Importance = widget.SuccessImportance
}

func getServicesState(_ fyne.Window) {
	ctx := context.Background()

	hlkLabel := aboutBodyContainer.Objects[0].(*fyne.Container).Objects[1].(*widget.Label)
	setServiceState(hlkLabel, hlkClient.GetIndex(ctx))

	hlpLabel := aboutBodyContainer.Objects[1].(*fyne.Container).Objects[1].(*widget.Label)
	setServiceState(hlpLabel, hlpClient.GetIndex(ctx))

	hlmLabel := aboutBodyContainer.Objects[2].(*fyne.Container).Objects[1].(*widget.Label)
	setServiceState(hlmLabel, hlmClient.GetIndex(ctx))

	hlfLabel := aboutBodyContainer.Objects[3].(*fyne.Container).Objects[1].(*widget.Label)
	setServiceState(hlfLabel, hlfClient.GetIndex(ctx))
}

func getFriends(w fyne.Window) {
	gFriends = []string{}
	friends, err := hlkClient.GetFriends(context.Background())
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	friendsList := proc.FriendsMapToList(friends)
	for _, f := range friendsList {
		gFriends = append(gFriends, f.FAliasName)
	}
}

func getConnectionSettings(w fyne.Window) {
	configSettings, err := hlkClient.GetSettings(context.Background())
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	text := encoding.SerializeJSON(configSettings)
	connectionSettingsLabel.SetText(string(text))
}

func getConnections(w fyne.Window) {
	gConnections = []sConnection{}

	connections, err := hlkClient.GetConnections(context.Background())
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	onlines, err := hlkClient.GetOnlines(context.Background())
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	mapOnlines := make(map[string]bool)
	for _, o := range onlines {
		mapOnlines[o] = true
	}

	for _, c := range connections {
		gConnections = append(gConnections, sConnection{
			online:  mapOnlines[c],
			address: c,
		})
	}
}

func pushMessage(w fyne.Window, friend, msg string) {
	t, err := hlmClient.PushMessage(context.Background(), friend, msg)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	addMessageToChat(w, friend, dto.NewMessage(false, msg, t))
	inputMessageEntry.SetText("")
}

func loadMessages(w fyne.Window, friend string) {
	msgs, err := hlmClient.LoadMessages(context.Background(), friend, 2048, 2048, true)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	for _, msg := range msgs {
		addMessageToChat(w, friend, msg)
	}
}

func listenMessages(w fyne.Window, friend string) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		chatListenerActive = true
		<-closeListenChat
		chatListenerActive = false
		cancel()
	}()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			msg, err := hlmClient.ListenChat(ctx, friend, "hidden-lake-client")
			if err != nil {
				select {
				case <-ctx.Done():
				case <-ticker.C:
				}
				continue
			}
			fyne.Do(func() { addMessageToChat(w, friend, msg) })
		}
	}()
}

func addMessageToChat(w fyne.Window, sender string, msg dto.IMessage) {
	var data fyne.CanvasObject
	msgData := msg.GetMessage()

	if strings.HasPrefix(msgData, "hls-filesharer:") {
		data = getMessageAsFile(w, sender, msg)
	} else {
		data = getMessageAsText(w, sender, msg)
	}

	c := container.New(
		layout.NewCustomPaddedVBoxLayout(0.1),
		func() *widget.Label {
			aliasName := sender
			if !msg.IsIncoming() {
				aliasName = "<IAM>"
			}
			msgLabel := widget.NewLabel(aliasName)
			msgLabel.Wrapping = fyne.TextWrapWord
			msgLabel.Selectable = true
			msgLabel.Importance = widget.HighImportance
			if msg.IsIncoming() {
				msgLabel.Importance = widget.DangerImportance
			}
			return msgLabel
		}(),
		data,
		func() *widget.Label {
			msgLabel := widget.NewLabel(msg.GetTimestamp())
			msgLabel.Wrapping = fyne.TextWrapWord
			msgLabel.Selectable = true
			msgLabel.Importance = widget.LowImportance
			return msgLabel
		}(),
	)

	bgColor := color.NRGBA{R: 0, G: 0, B: 0, A: 128}
	backgroundRect := canvas.NewRectangle(bgColor)
	coloredContainer := container.NewStack(backgroundRect, c)

	scrollChatContainer.Content.(*fyne.Container).Add(coloredContainer)
	scrollChatContainer.ScrollToBottom()
}

func getMessageAsText(_ fyne.Window, _ string, msg dto.IMessage) *widget.Label {
	msgLabel := widget.NewLabel(msg.GetMessage())
	msgLabel.Wrapping = fyne.TextWrapWord
	msgLabel.Selectable = true
	return msgLabel
}

func getMessageAsFile(w fyne.Window, sender string, msg dto.IMessage) *fyne.Container {
	msgData := msg.GetMessage()

	filename := strings.Replace(msgData, "hls-filesharer:", "", 1)
	viewFilename := strings.Replace(filename, filepath.Ext(filename), "", 1)

	downloadButton := widget.NewButtonWithIcon("LOAD", theme.DownloadIcon(), func() {
		fileDialog := dialog.NewFileSave(
			func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if writer == nil {
					return
				}

				go func() {
					defer writer.Close()

					if msg.IsIncoming() {
						downloading, err := hlfClient.GetRemoteFile(writer, context.Background(), sender, filename, true)
						if err != nil {
							fyne.Do(func() { dialog.ShowError(err, w) })
							return
						}
						if downloading {
							fyne.Do(func() { dialog.ShowInformation("Download state", "File downloading...", w) })
							return
						}
					} else {
						if err := hlfClient.GetLocalFile(writer, context.Background(), sender, filename); err != nil {
							fyne.Do(func() { dialog.ShowError(err, w) })
							return
						}
					}
					fyne.Do(func() { dialog.ShowInformation("Download state", "File success downloaded!", w) })
				}()
			},
			w,
		)
		fileDialog.SetFileName(viewFilename)
		fileDialog.Show()
	})
	downloadButton.Importance = widget.LowImportance

	infoButton := widget.NewButtonWithIcon("INFO", theme.FileIcon(), func() {
		go func() {
			var (
				fileInfo dto_f.IFileInfo
				err      error
			)
			if msg.IsIncoming() {
				fileInfo, err = hlfClient.GetRemoteFileInfo(context.Background(), sender, filename, true)
			} else {
				fileInfo, err = hlfClient.GetLocalFileInfo(context.Background(), sender, filename)
			}
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			fyne.Do(func() { dialog.ShowInformation("File info", fileInfo.ToString(), w) })
		}()
	})
	infoButton.Importance = widget.LowImportance

	buttons := container.NewGridWithColumns(
		2,
		downloadButton,
		infoButton,
	)

	msgLabel := widget.NewLabel(viewFilename)
	msgLabel.Importance = widget.WarningImportance
	msgLabel.Wrapping = fyne.TextWrapWord
	msgLabel.Selectable = true

	return container.New(
		layout.NewBorderLayout(nil, nil, nil, buttons),
		msgLabel,
		buttons,
	)
}
