package model

// Todo -
type Todo struct {
	ID     int    `json:"id" gorm:"primary_key,auto_increment"`
	Text   string `json:"text" gorm:"type:text"`
	Done   bool   `json:"done"`
	UserID int    `gorm:"index"`

	// -----------
	BaseModel

	// -------
	User *User `json:"user" gorm:"-"`
}
