package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	aboutContainer       = new(fyne.Container)
	chatListContainer    = new(fyne.Container)
	chatFriendContainer  = new(fyne.Container)
	friendsContainer     = new(fyne.Container)
	connectionsContainer = new(fyne.Container)
)

func main() {
	a := app.NewWithID("hidden.lake.client")

	w := a.NewWindow("Hidden Lake Client")
	w.Resize(fyne.NewSize(600, 400))

	aboutContainer = initWindowAbout(a, w)
	chatListContainer = initWindowChatList(a, w)
	chatFriendContainer = initWindowChatFriend(a, w)
	friendsContainer = initWindowFriends(a, w)
	connectionsContainer = initWindowConnections(a, w)

	setChatListContent(w)
	w.ShowAndRun()
}
