package message

type IMessageContainer interface {
	GetFriend() string
	GetMessageBody() []byte
}
