package storage

type Storage interface {
	AddChatId(int64) error
	RemoveChatId(int64) error
	LoadAllChatId() ([]int64, error)
	GetAuthToken() (string, error)
}
