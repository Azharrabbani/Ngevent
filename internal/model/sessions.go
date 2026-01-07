package model

type Sessions struct {
	ID           string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	User         Users  `gorm:"foreignKey:UserID"`
}

func (Sessions) TableName() string {
	return "sessions"
}
