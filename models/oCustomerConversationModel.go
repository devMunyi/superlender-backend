package models

import (
	"time"
)

type ConversationStatus uint8

const (
	DeletedConversation ConversationStatus = iota
	ActiveConversation  ConversationStatus = 1
	BlockedConversation ConversationStatus = 2
)

type OCustomerConversation struct {
	UID                    int                `json:"uid" gorm:"primaryKey;autoIncrement"`
	CustomerID             int                `json:"customer_id" gorm:"not null"`
	Branch                 int                `json:"branch" gorm:"not null"`
	AgentID                int                `json:"agent_id" gorm:"default:0"`
	LoanID                 int                `json:"loan_id" gorm:"default:0"`
	Transcript             string             `json:"transcript" gorm:"type:mediumtext;not null"`
	ConversationMethod     int                `json:"conversation_method" gorm:"not null"`
	ConversationDate       time.Time          `json:"conversation_date" gorm:"autoCreateTime;type:datetime"`
	ConversationDateNoTime time.Time          `json:"conversation_date_without_time" gorm:"autoCreateTime;type:date"`
	NextInteraction        time.Time          `json:"next_interaction" gorm:"type:date;not null"`
	NextSteps              int                `json:"next_steps" gorm:"not null"`
	Flag                   int                `json:"flag" gorm:"not null"`
	Outcome                int                `json:"outcome" gorm:"default:0"`
	Status                 ConversationStatus `json:"status" gotm:"default:1"`
}
