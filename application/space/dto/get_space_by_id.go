package dto

import "time"

// SpacePermission represents a permission entry for a space
type SpacePermission struct {
	UserId  *string
	GroupId *string
	Role    string
}

// SpaceDetail contains the space data needed by other applications
type SpaceDetail struct {
	Id          string
	Name        string
	Slug        string
	Icon        string
	IconColor   string
	Type        string
	OwnerId     *string
	Permissions []SpacePermission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetUserRole returns the user's role in this space, or nil if no role
func (s *SpaceDetail) GetUserRole(userId string) *string {
	// Check if the user is the owner
	if s.OwnerId != nil && *s.OwnerId == userId {
		role := "owner"
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

// HasPermission checks if the user has at least the required role level
func (s *SpaceDetail) HasPermission(userId string, requiredRole string) bool {
	userRole := s.GetUserRole(userId)
	if userRole == nil {
		// For public spaces, allow reading
		return s.Type == "public" && requiredRole == "viewer"
	}

	return roleHasPermission(*userRole, requiredRole)
}

func roleHasPermission(userRole, requiredRole string) bool {
	roleHierarchy := map[string]int{
		"viewer": 1,
		"editor": 2,
		"admin":  3,
		"owner":  4,
	}
	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}

type GetSpaceByIdInput struct {
	SpaceId string
}

type GetSpaceByIdOutput struct {
	Space *SpaceDetail
}
