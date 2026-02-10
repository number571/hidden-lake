package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
)

func initWindowAbout(a fyne.App, w fyne.Window) *fyne.Container {
	header := widget.NewButtonWithIcon(
		"Back to main page",
		theme.ListIcon(),
		func() { setChatListContent(w) },
	)

	versionGrid := container.NewGridWithColumns(
		2,
		widget.NewLabel("Version"),
		widget.NewLabel(version),
	)
	versionGrid.Objects[1].(*widget.Label).Importance = widget.WarningImportance

	coloredVersionGridContainer := container.NewStack(
		canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 100}),
		versionGrid,
	)

	aboutBodyContainer = container.NewGridWithRows(
		4,
		container.NewGridWithColumns(
			2,
			widget.NewLabel("HLK"),
			widget.NewLabel(""),
		),
		container.NewGridWithColumns(
			2,
			widget.NewLabel("HLS=pinger"),
			widget.NewLabel(""),
		),
		container.NewGridWithColumns(
			2,
			widget.NewLabel("HLS=messenger"),
			widget.NewLabel(""),
		),
		container.NewGridWithColumns(
			2,
			widget.NewLabel("HLS=filesharer"),
			widget.NewLabel(""),
		),
	)

	coloredAboutBodyLabelContainer := container.NewStack(
		canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 100}),
		aboutBodyContainer,
	)

	hyperlinkToAuthorWithLabel := container.NewGridWithColumns(
		2,
		widget.NewLabel("Author:"),
		widget.NewHyperlink("github.com/number571", func() *url.URL {
			urlObj, _ := url.Parse("https://github.com/number571")
			return urlObj
		}()),
	)

	hyperlinkToProjectWithLabel := container.NewGridWithColumns(
		2,
		widget.NewLabel("Project:"),
		widget.NewHyperlink("github.com/number571/hidden-lake", func() *url.URL {
			urlObj, _ := url.Parse("https://github.com/number571/hidden-lake")
			return urlObj
		}()),
	)

	hyperlinkToWhitePaperWithLabel := container.NewGridWithColumns(
		2,
		widget.NewLabel("White paper:"),
		widget.NewHyperlink("hidden_lake_anonymous_network.pdf", func() *url.URL {
			urlObj, _ := url.Parse("https://github.com/number571/hidden-lake/blob/master/docs/hidden_lake_anonymous_network.pdf")
			return urlObj
		}()),
	)

	gridOfHyperlinks := container.NewGridWithRows(
		3,
		hyperlinkToAuthorWithLabel,
		hyperlinkToProjectWithLabel,
		hyperlinkToWhitePaperWithLabel,
	)

	coloredHyperlinkWithLabels := container.NewStack(
		canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 100}),
		gridOfHyperlinks,
	)

	innerContent := container.NewVBox(
		coloredVersionGridContainer,
		coloredAboutBodyLabelContainer,
		coloredHyperlinkWithLabels,
	)

	content := container.New(
		layout.NewBorderLayout(header, nil, nil, nil),
		header,
		innerContent,
	)

	minSizeTarget := canvas.NewRectangle(color.Transparent)
	minSizeTarget.SetMinSize(fyne.NewSize(600, 400))

	contentContainerWrapper := container.New(
		layout.NewStackLayout(),
		minSizeTarget,
		content,
	)

	w.SetCloseIntercept(func() { a.Quit() })
	return contentContainerWrapper
}

func initWindowFriends(a fyne.App, w fyne.Window) *fyne.Container {
	header := widget.NewButtonWithIcon(
		"Back to main page",
		theme.ListIcon(),
		func() { setChatListContent(w) },
	)

	pubKeyButton := widget.NewButtonWithIcon(
		"Copy my public key",
		theme.ContentCopyIcon(),
		func() {
			pubKey, err := hlkClient.GetPubKey(context.Background())
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			a.Clipboard().SetContent(pubKey.ToString())
			dialog.ShowInformation(
				"Copying a public key...",
				"The public key has been successfully copied to the clipboard",
				w,
			)
		},
	)
	pubKeyButton.Importance = widget.LowImportance

	coloredPubKeyButtonContainer := container.NewStack(
		canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 100}),
		pubKeyButton,
	)

	inputFriendNameEntry = widget.NewEntry()
	inputFriendNameEntry.SetPlaceHolder("Type a name...")

	inputFriendPubKeyEntry = widget.NewEntry()
	inputFriendPubKeyEntry.SetPlaceHolder("Type a key...")

	sendButton := widget.NewButtonWithIcon("", theme.MailForwardIcon(), func() {
		aliasName := inputFriendNameEntry.Text
		if aliasName == "" {
			dialog.ShowError(errors.New("invalid alias name"), w)
			return
		}
		pubKey := asymmetric.LoadPubKey(inputFriendPubKeyEntry.Text)
		if pubKey == nil {
			dialog.ShowError(errors.New("invalid public key"), w)
			return
		}
		if err := hlkClient.AddFriend(context.Background(), aliasName, pubKey); err != nil {
			dialog.ShowError(err, w)
			return
		}
		setFriendsContent(w)
	})

	inputFriendNameEntry.OnSubmitted = func(s string) {
		sendButton.Tapped(nil)
	}
	inputFriendPubKeyEntry.OnSubmitted = func(s string) {
		sendButton.Tapped(nil)
	}

	friendsList := widget.NewList(
		func() int { return len(gFriends) },
		func() fyne.CanvasObject {
			templateDeleteButton := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {})
			templateDeleteButton.Importance = widget.DangerImportance

			templatePingButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameDragCornerIndicator), func() {})
			templateFriendButton := widget.NewButtonWithIcon("", theme.AccountIcon(), func() {})

			return container.New(
				layout.NewBorderLayout(nil, nil, templateDeleteButton, templatePingButton),
				templateDeleteButton,
				templatePingButton,
				templateFriendButton,
			)
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			friend := gFriends[i]

			deleteButton := item.(*fyne.Container).Objects[0].(*widget.Button)
			deleteButton.OnTapped = func() {
				dialog.ShowConfirm(
					"Deleting friend...",
					"Are you sure you want to delete this friend?",
					func(ok bool) {
						if !ok {
							return
						}
						if err := hlkClient.DelFriend(context.Background(), friend); err != nil {
							dialog.ShowError(err, w)
						}
						setFriendsContent(w)
					},
					w,
				)
			}

			pingButton := item.(*fyne.Container).Objects[1].(*widget.Button)
			pingButton.OnTapped = func() {
				go func() {
					defer fyne.Do(func() { pingButton.Refresh() })
					if err := hlpClient.PingFriend(context.Background(), friend); err != nil {
						pingButton.Importance = widget.WarningImportance
						return
					}
					pingButton.Importance = widget.SuccessImportance
				}()
			}

			friendButton := item.(*fyne.Container).Objects[2].(*widget.Button)
			friendButton.SetText(friend)
			friendButton.OnTapped = func() {
				friendsMap, err := hlkClient.GetFriends(context.Background())
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				pubKey, ok := friendsMap[friend]
				if !ok {
					dialog.ShowError(errors.New("not found friend"), w)
					return
				}
				a.Clipboard().SetContent(pubKey.ToString())
				dialog.ShowInformation(
					"Copying a friends public key...",
					"The public key has been successfully copied to the clipboard",
					w,
				)
			}
		},
	)

	inputBarInner := container.NewGridWithColumns(
		3,
		inputFriendNameEntry,
		inputFriendPubKeyEntry,
		sendButton,
	)

	inputBar := container.New(
		layout.NewBorderLayout(coloredPubKeyButtonContainer, nil, nil, nil),
		coloredPubKeyButtonContainer,
		inputBarInner,
	)

	content := container.New(
		layout.NewBorderLayout(header, inputBar, nil, nil),
		header,
		friendsList,
		inputBar,
	)

	minSizeTarget := canvas.NewRectangle(color.Transparent)
	minSizeTarget.SetMinSize(fyne.NewSize(600, 400))

	contentContainerWrapper := container.New(
		layout.NewStackLayout(),
		minSizeTarget,
		content,
	)

	w.SetCloseIntercept(func() { a.Quit() })
	return contentContainerWrapper
}

func initWindowConnections(a fyne.App, w fyne.Window) *fyne.Container {
	header := widget.NewButtonWithIcon(
		"Back to main page",
		theme.ListIcon(),
		func() { setChatListContent(w) },
	)

	networksList := widget.NewList(
		func() int { return len(gConnections) },
		func() fyne.CanvasObject {
			templateNetworkButton := widget.NewButtonWithIcon("", theme.ListIcon(), func() {})

			templateDeleteNetwork := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {})
			templateDeleteNetwork.Importance = widget.DangerImportance

			return container.New(
				layout.NewBorderLayout(nil, nil, templateDeleteNetwork, nil),
				templateDeleteNetwork,
				templateNetworkButton,
			)
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			connection := gConnections[i]

			deleteButton := item.(*fyne.Container).Objects[0].(*widget.Button)
			deleteButton.OnTapped = func() {
				dialog.ShowConfirm(
					"Deleting connection...",
					"Are you sure you want to delete this connection?",
					func(ok bool) {
						if !ok {
							return
						}
						if err := hlkClient.DelConnection(context.Background(), connection.address); err != nil {
							dialog.ShowError(err, w)
							return
						}
						setConnectionsContent(w)
					},
					w,
				)
			}

			buttonName := item.(*fyne.Container).Objects[1].(*widget.Button)
			buttonName.SetText(connection.address)
			if connection.online {
				buttonName.Importance = widget.SuccessImportance
			}

			buttonName.OnTapped = func() {
				a.Clipboard().SetContent(connection.address)
				dialog.ShowInformation(
					"Copying a connection...",
					"The connection has been successfully copied to the clipboard",
					w,
				)
			}
		},
	)

	inputConnectionEntry = widget.NewEntry()
	inputConnectionEntry.SetPlaceHolder("Type a connection...")

	sendButton := widget.NewButtonWithIcon(
		"",
		theme.MailForwardIcon(),
		func() {
			connection := inputConnectionEntry.Text
			inputConnectionEntry.SetText("")
			if err := hlkClient.AddConnection(context.Background(), connection); err != nil {
				dialog.ShowError(err, w)
				return
			}
			setConnectionsContent(w)
		},
	)

	inputConnectionEntry.OnSubmitted = func(s string) {
		sendButton.Tapped(nil)
	}

	connectionSettingsLabel = widget.NewLabel("")
	connectionSettingsLabel.Selectable = true
	connectionSettingsLabel.Wrapping = fyne.TextWrapWord

	coloredLabelContainer := container.NewStack(
		canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 100}),
		connectionSettingsLabel,
	)

	scrollContainer := container.NewScroll(coloredLabelContainer)
	scrollContainer.SetMinSize(fyne.NewSize(600, 100))

	inputEntrySendButton := container.New(
		layout.NewBorderLayout(nil, nil, nil, sendButton),
		inputConnectionEntry,
		sendButton,
	)

	inputEntrySendButtonWithSettingsInfo := container.New(
		layout.NewBorderLayout(nil, scrollContainer, nil, nil),
		networksList,
		scrollContainer,
	)

	content := container.New(
		layout.NewBorderLayout(header, inputEntrySendButton, nil, nil),
		header,
		inputEntrySendButtonWithSettingsInfo,
		inputEntrySendButton,
	)

	minSizeTarget := canvas.NewRectangle(color.Transparent)
	minSizeTarget.SetMinSize(fyne.NewSize(600, 400))

	contentContainerWrapper := container.New(
		layout.NewStackLayout(),
		minSizeTarget,
		content,
	)

	w.SetCloseIntercept(func() { a.Quit() })
	return contentContainerWrapper
}

func initWindowChatFriend(a fyne.App, w fyne.Window) *fyne.Container {
	scrollChatContainer = container.NewScroll(container.NewVBox())
	scrollChatContainer.SetMinSize(fyne.NewSize(400, 300))

	inputMessageEntry = widget.NewEntry()
	inputMessageEntry.SetPlaceHolder("Type a message...")

	fileButton := widget.NewButtonWithIcon("", theme.FileIcon(), func() {
		fileOpenDialog := dialog.NewFileOpen(
			func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if reader == nil {
					return
				}
				defer reader.Close()

				filename := fmt.Sprintf("%s.%s", filepath.Base(reader.URI().Name()), random.NewRandom().GetString(16))
				if err := hlfClient.PutLocalFile(context.Background(), friendNameInChat, filename, reader); err != nil {
					dialog.ShowError(err, w)
					return
				}
				pushMessage(w, friendNameInChat, fmt.Sprintf("hls-filesharer:%s", filename))
			},
			w,
		)
		fileOpenDialog.Show()
	})

	sendButton := widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		content := inputMessageEntry.Text
		if content == "" {
			return
		}
		pushMessage(w, friendNameInChat, content)
	})

	sendButtons := container.New(
		layout.NewBorderLayout(nil, nil, nil, sendButton),
		fileButton,
		sendButton,
	)

	inputBar := container.New(
		layout.NewBorderLayout(nil, nil, nil, sendButtons),
		inputMessageEntry,
		sendButtons,
	)

	inputMessageEntry.OnSubmitted = func(s string) {
		sendButton.Tapped(nil)
	}

	header := widget.NewButtonWithIcon(
		"Back to main page",
		theme.ListIcon(),
		func() { setChatListContent(w) },
	)

	content := container.New(
		layout.NewBorderLayout(header, inputBar, nil, nil),
		header,
		inputBar,
		scrollChatContainer,
	)

	minSizeTarget := canvas.NewRectangle(color.Transparent)
	minSizeTarget.SetMinSize(fyne.NewSize(600, 400))

	contentContainerWrapper := container.New(
		layout.NewStackLayout(),
		minSizeTarget,
		content,
	)

	w.SetCloseIntercept(func() { a.Quit() })
	return contentContainerWrapper
}

func initWindowChatList(a fyne.App, w fyne.Window) *fyne.Container {
	chatList := widget.NewList(
		func() int {
			return len(gFriends)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewButton("", func() {}))
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			friend := gFriends[i]

			buttonName := item.(*fyne.Container).Objects[0].(*widget.Button)
			buttonName.SetText(friend)
			buttonName.OnTapped = func() { setChatFriendContent(w, friend) }
		},
	)

	mainContentVBox := container.NewBorder(nil, nil, nil, nil, chatList)
	networksButton := widget.NewButtonWithIcon(
		"",
		theme.SettingsIcon(),
		func() { setConnectionsContent(w) },
	)
	aboutButton := widget.NewButtonWithIcon(
		"",
		theme.MenuIcon(),
		func() { setAboutContent(w) },
	)
	friendButton := widget.NewButtonWithIcon(
		"Friends",
		theme.ComputerIcon(),
		func() { setFriendsContent(w) },
	)

	header := container.New(
		layout.NewBorderLayout(nil, nil, networksButton, aboutButton),
		friendButton,
		networksButton,
		aboutButton,
	)

	content := container.New(
		layout.NewBorderLayout(header, nil, nil, nil),
		header,
		mainContentVBox,
	)

	minSizeTarget := canvas.NewRectangle(color.Transparent)
	minSizeTarget.SetMinSize(fyne.NewSize(600, 400))

	contentContainerWrapper := container.New(
		layout.NewStackLayout(),
		minSizeTarget,
		content,
	)

	w.SetCloseIntercept(func() { a.Quit() })
	return contentContainerWrapper
}
