package domain

import (
	"time"

	"gorm.io/gorm"
)

type Space struct {
	Id   string
	Name string

	Slug      string
	Icon      string
	IconColor string

	Type SpaceType

	OwnerId *string
	Owner   *User `gorm:"foreignKey:OwnerId;references:Id"`

	Documents []Document `gorm:"foreignKey:SpaceId;references:Id"`

	Permissions []Permission `gorm:"foreignKey:SpaceId;references:Id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// SpaceType is the type of space
type SpaceType string

const (
	// SpaceTypePublic is the public space type
	SpaceTypePublic SpaceType = "public"
	// SpaceTypePrivate is the private space type
	SpaceTypePrivate SpaceType = "private"
	// SpaceTypeRestricted is the restricted space type
	SpaceTypeRestricted SpaceType = "restricted"
	// SpaceTypePersonal is the personal space type
	SpaceTypePersonal SpaceType = "personal"
)

func (s *Space) TableName() string {
	return "space"
}

func (s *Space) GetUserRole(userId string) *PermissionRole {
	// Check if the user is the owner
	if s.OwnerId != nil && *s.OwnerId == userId {
		role := PermissionRoleOwner
		return &role
	}

	// Check permissions
	for _, perm := range s.Permissions {
		if perm.UserId != nil && *perm.UserId == userId {
			return &perm.Role
		}
	}

	return nil
}

func (s *Space) HasPermission(userId string, requiredRole PermissionRole) bool {
	userRole := s.GetUserRole(userId)
	if userRole == nil {
		// For public spaces, allow reading
		return s.Type == SpaceTypePublic && requiredRole == PermissionRoleViewer
	}

	return s.roleHasPermission(*userRole, requiredRole)
}

func (s *Space) roleHasPermission(userRole, requiredRole PermissionRole) bool {
	roleHierarchy := map[PermissionRole]int{
		PermissionRoleViewer: 1,
		PermissionRoleEditor: 2,
		PermissionRoleAdmin:  3,
		PermissionRoleOwner:  4,
	}
	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}

type SpacePers interface {
	Create(space *Space) error
	GetSpacesForUser(userId string) ([]Space, error)
	GetSpaceById(spaceId string) (*Space, error)
	Update(space *Space) error
	Delete(spaceId string) error
	// Admin methods
	GetAll(limit, offset int) ([]Space, int64, error)
}
