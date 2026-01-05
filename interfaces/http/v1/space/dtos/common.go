package dtos

import "time"

type Space struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Icon      string `json:"icon"`
	IconColor string `json:"icon_color"`
	Type      string `json:"type"`
	MyRole    string `json:"my_role,omitempty"` // Role of the current user in this space (owner, admin, editor, viewer)

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
