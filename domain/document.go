package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Document struct {
	Id   string
	Name string
	Slug string

	Config   DocumentConfig
	Metadata JSONB

	ParentId *string
	Parent   *Document `gorm:"foreignKey:ParentId;references:Id"`

	SpaceId string
	Space   Space `gorm:"foreignKey:SpaceId;references:Id"`

	Public bool

	// Permissions spécifiques au document (optionnelles)
	Permissions []DocumentPermission `gorm:"foreignKey:DocumentId;references:Id"`

	Content datatypes.JSON

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (d *Document) TableName() string {
	return "document"
}

type DocumentConfig struct {
	FullWidth        bool   `json:"full_width"`
	Icon             string `json:"icon"`
	Lock             bool   `json:"lock"`
	HeaderBackground string `json:"header_background"`
}

// Value implements the driver.Valuer interface
func (dc DocumentConfig) Value() (driver.Value, error) {
	return json.Marshal(dc)
}

// Scan implements the sql.Scanner interface
func (dc *DocumentConfig) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		// PostgreSQL usually returns []byte
		return json.Unmarshal(v, dc)
	case string:
		// SQLite often returns string
		return json.Unmarshal([]byte(v), dc)
	case nil:
		// Handle null case
		*dc = DocumentConfig{}
		return nil
	default:
		// Fall back to string conversion
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, dc)
	}
}

func (d *Document) HasPermission(userId string, requiredRole DocumentRole) bool {
	// 1. Vérifier les permissions spécifiques au document d'abord
	for _, perm := range d.Permissions {
		if perm.UserId != nil && *perm.UserId == userId {
			if perm.Role == DocumentRoleDenied {
				return false // Refus explicite
			}
			return d.documentRoleHasPermission(perm.Role, requiredRole)
		}
	}

	// 2. Si pas de permission spécifique, hériter du space
	if d.Space.HasPermission(userId, SpaceRoleViewer) {
		// Si l'user a accès au space, il peut au moins voir le document
		if requiredRole == DocumentRoleViewer {
			return true
		}
		// Pour éditer, il faut au moins être editor du space
		if requiredRole == DocumentRoleEditor {
			return d.Space.HasPermission(userId, SpaceRoleEditor)
		}
	}

	return false
}

func (d *Document) documentRoleHasPermission(userRole, requiredRole DocumentRole) bool {
	roleHierarchy := map[DocumentRole]int{
		DocumentRoleViewer: 1,
		DocumentRoleEditor: 2,
		DocumentRoleOwner:  3,
	}
	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}

// CanManagePermissions returns true if the user can manage document permissions
// This requires being owner of the document OR admin/owner of the space
func (d *Document) CanManagePermissions(userId string) bool {
	// Check if user is owner of this document
	for _, perm := range d.Permissions {
		if perm.UserId != nil && *perm.UserId == userId && perm.Role == DocumentRoleOwner {
			return true
		}
	}

	// Check if user is admin or owner of the space
	return d.Space.HasPermission(userId, SpaceRoleAdmin)
}

type DocumentPers interface {
	GetDocumentWithPermissions(documentId, userId string) (*Document, error)
	GetDocumentByIdOrSlugWithUserPermissions(spaceId string, id *string, slug *string, userId string) (*Document, error)
	GetRootDocumentsFromSpaceWithUserPermissions(spaceId, userId string) ([]Document, error)
	GetChildDocumentsWithUserPermissions(parentId, userId string) ([]Document, error)
	Create(document *Document, userId string) error
	Update(document *Document, userId string) error
	Delete(documentId, userId string) error
	Move(documentId string, newParentId *string, userId string) (*Document, error)
	// Trash management
	GetDeletedDocuments(spaceId, userId string) ([]Document, error)
	Restore(documentId, userId string) error
	// Public sharing
	SetPublic(documentId string, public bool, userId string) error
	GetPublicDocument(spaceId string, id *string, slug *string) (*Document, error)
	// Search
	Search(query string, userId string, spaceId *string, limit int) ([]Document, error)
}
