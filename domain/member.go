package domain

type Members []Member
type MembersWithUsersOrGroups []MemberWithUsersOrGroups

type MemberType string
type AccessType string

type Member struct {
	Id     string
	Type   MemberType
	Access AccessType
}

// AccessTypeViewer is the access type viewer
const (
	AccessTypeViewer  AccessType = "viewer"
	AccessTypeEditor  AccessType = "editor"
	AccessTypeComment AccessType = "comment"
	AccessTypeFull    AccessType = "full"
)

// MemberType is the type of member
const (
	MemberTypeUser  MemberType = "user"
	MemberTypeGroup MemberType = "group"
)

// MemberWithUser is a model for a member with user information
type MemberWithUsersOrGroups struct {
	Member
	User  User
	Group Group
}
