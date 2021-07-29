package api

type CreateGroupRequest struct {
	Owner      string   `form:"owner" json:"owner"`
	Type       string   `form:"type" json:"type" binding:"required"`
	GroupId    string   `form:"group_id" json:"group_id"`
	Name       string   `form:"name" json:"name" binding:"required"`
	MemberList []string `form:"member_list" json:"member_list"`
}

type CreateGroupResponse struct {
	GroupId string `json:"group_id,omitempty"`
}

type GroupProfile struct {
	// 群资料
}
