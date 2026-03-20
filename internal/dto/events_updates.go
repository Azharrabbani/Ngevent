package dto

type EventsUpdatesResp struct {
	ID                string            `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID           string            `json:"event_id"`
	EventTitle        string            `json:"event_title"`
	UpdatedDetails    UpdatedDetails    `json:"updated_details"`
	UpdatedAddress    UpdatedAddress    `json:"updated_address"`
	UpdatedCategories []EventCategories `json:"updated_categories"`
	UpdatedTickets    int               `json:"updated_tickets"`
	CreatedAt         int64             `json:"created_at" gorm:"default:now()"`
	UpdatedAt         int64             `json:"updated_at" gorm:"default:now()"`
	DeletedAt         int64             `json:"deleted_at"`
}

type UpdatedDetails struct {
	Banner      string `json:"banner,omitempty"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Date        int64  `json:"date"`
}

type UpdatedAddress struct {
	Address       string `json:"address"`
	City          string `json:"city"`
	Country       string `json:"country"`
	DetailAddress string `json:"detail_address"`
	Coordinates   string `json:"cooridinates"`
}
