package models

type ConversationMethodStatus uint8

const (
	DeletedConversationMethod ConversationMethodStatus = iota
	ActiveConversationMethod  ConversationMethodStatus = 1
)

type OConversationMethod struct {
	UID     int                      `json:"uid" gorm:"primaryKey;autoIncrement"`
	Name    string                   `json:"name" gorm:"unique;not null"`
	Details string                   `json:"details" gorm:"not null"`
	Status  ConversationMethodStatus `json:"status" gorm:"default:1"`
}
