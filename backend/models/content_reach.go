package models

type ContentReach string

const (
	ContentReachPublic  ContentReach = "public"
	ContentReachMembers ContentReach = "members"
	ContentReachStaff   ContentReach = "staff"
)

var ValidContentReach = []ContentReach{
	ContentReachPublic,
	ContentReachMembers,
	ContentReachStaff,
}

func DefaultContentReach(visibility EventVisibility) ContentReach {
	if visibility == EventVisibilityPublic {
		return ContentReachPublic
	}
	return ContentReachMembers
}
