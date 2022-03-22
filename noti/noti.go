package noti

type TelegramNoti interface {
	Send(string) error
}
