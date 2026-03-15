package dto

// DocumentPermission represents a permission entry for a document
type DocumentPermission struct {
	UserId *string
	Role   string
}

// DocumentSpaceInfo contains the space info embedded in a document
type DocumentSpaceInfo struct {
	Type    string
	OwnerId *string
}

// DocumentDetail contains the document data needed by other applications
type DocumentDetail struct {
	Id          string
	Name        string
	Slug        string
	SpaceId     string
	Config      DocumentConfig
	Space       DocumentSpaceInfo
	Permissions []DocumentPermission
}

// HasPermission checks if a user has at least the required role on this document
func (d *DocumentDetail) HasPermission(userId string, requiredRole string) bool {
	// 1. Check document-level permissions first
	for _, perm := range d.Permissions {
		if perm.UserId != nil && *perm.UserId == userId {
			if perm.Role == "denied" {
				return false
			}
			return docRoleHasPermission(perm.Role, requiredRole)
		}
	}

	// 2. Inherit from space - simplified check
	// For public spaces, viewers can see
	if d.Space.Type == "public" && requiredRole == "viewer" {
		return true
	}

	return false
}

// CanManagePermissions returns true if the user can manage document permissions
func (d *DocumentDetail) CanManagePermissions(userId string) bool {
	// Check if user is owner of this document
	for _, perm := range d.Permissions {
		if perm.UserId != nil && *perm.UserId == userId && perm.Role == "owner" {
			return true
		}
	}

	// Space admin/owner check would need the space permissions loaded
	// The space-level check is done via the space data embedded in the document
	return false
}

func docRoleHasPermission(userRole, requiredRole string) bool {
	roleHierarchy := map[string]int{
		"viewer": 1,
		"editor": 2,
		"owner":  3,
	}
	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}

type GetDocumentByIdOrSlugWithUserPermissionsInput struct {
	DocumentId *string
	Slug       *string
	SpaceId    string
	UserId     string
}

type GetDocumentByIdOrSlugWithUserPermissionsOutput struct {
	Document *DocumentDetail
}
