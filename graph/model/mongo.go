package model

type Note struct {
	ID        string `json:"_id" bson:"_id,omitempty"`
	Title     string `json:"title" bson:"title,omitempty"`
	Content   string `json:"content" bson:"content,omitempty"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt int64  `json:"updatedAt" bson:"updatedAt,omitempty"`
	CreatedBy *User  `json:"createdBy" bson:"createdBy,omitempty"`
}

type User struct {
	ID        string  `json:"_id" bson:"_id,omitempty"`
	Name      string  `json:"name" bson:"name,omitempty"`
	Email     string  `json:"email" bson:"email,omitempty"`
	Password  string  `json:"password" bson:"password,omitempty"`
	CreatedAt int64   `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt int64   `json:"updatedAt" bson:"updatedAt,omitempty"`
	Notes     []*Note `json:"notes" bson:"notes,omitempty"`
}
