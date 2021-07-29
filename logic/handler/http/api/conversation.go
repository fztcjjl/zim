package api

type GetRecentConversationRequest struct {
	UserId string `form:"user_id" json:"user_id"`
	Offset int64  `form:"offset" json:"offset"`
	Limit  int64  `form:"limit" json:"limit"`
}

type Conversation struct {
	Type            int    `json:"type"`
	Target          string `json:"target"`
	UnreadCount     int    `json:"unread_count"`
	C2CPeerReadTime int64  `json:"c2c_peer_read_time"`
	IsTop           bool   `json:"is_top"`
	IsMute          bool   `json:"is_mute"`
	MsgShow         string `json:"msg_show"`
	Remark          string `json:"remark"`
	LastMsg         *Msg   `json:"last_msg"`
}

type GetRecentConversationResponse struct {
	List []*Conversation `json:"list"`
}

type GetConversationMsgRequest struct {
	ConvId string `form:"conv_id" json:"conv_id"`
	Offset int64  `form:"offset" json:"offset"`
	Limit  int64  `form:"limit" json:"limit"`
}

type GetConversationMsgResponse struct {
	List []*Msg `json:"list"`
}

type GetConversationRequest struct {
	ConvId string `form:"conv_id" json:"conv_id"`
}

type GetConversationResponse struct {
	Conversation
}

type SetConversationTopRequest struct {
	ConvId string `form:"conv_id" json:"conv_id"`
	IsTop  bool   `form:"is_top" json:"is_top"`
}

type SetConversationTopResponse struct {
}

type SetConversationMuteRequest struct {
	ConvId string `form:"conv_id" json:"conv_id"`
	IsMute bool   `form:"is_mute" json:"is_mute"`
}

type SetConversationMuteResponse struct {
}

type SetConversationReadRequest struct {
	ConvId string `form:"conv_id" json:"conv_id"`
}

type SetConversationReadResponse struct {
}

type GetPeerReadTimeRequest struct {
	ConvId string `form:"conv_id" json:"conv_id"`
}

type GetPeerReadTimeResponse struct {
	C2CPeerReadTime int64 `form:"c2c_peer_read_time" json:"c2c_peer_read_time"`
}
