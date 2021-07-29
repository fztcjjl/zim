package api

type Msg struct {
	Id            int64    `json:"id"`
	ConvType      int      `json:"conv_type"`
	Type          int      `json:"type"`
	Content       string   `json:"content"`
	Sender        string   `json:"sender"`
	Target        string   `json:"target"`
	Extra         string   `json:"extra"`
	SendTime      int64    `json:"send_time"`
	AtUserList    []string `json:"at_user_list"`
	IsTransparent bool     `json:"is_transparent"`
}

type TextMsg struct {
	Text string `json:"text"`
}

type SendRequest struct {
	ConvType      int      `form:"conv_type" json:"conv_type" binding:"required"`
	MsgType       int      `form:"msg_type" json:"msg_type" binding:"required"`
	Sender        string   `form:"sender" json:"sender" binding:"required"`
	Target        string   `form:"target" json:"target" binding:"required"`
	Content       string   `form:"content" json:"content" binding:"required"`
	Extra         string   `form:"extra" json:"extra"`
	AtUserList    []string `form:"at_user_list" json:"at_user_list"`
	IsTransparent bool     `form:"is_transparent" json:"is_transparent"`
}

type SendResponse struct {
	Id       int64 `json:"id,omitempty"`
	SendTime int64 `json:"send_time,omitempty"`
	Seq      int64 `json:"seq,omitempty"`
}

type SyncMsgRequest struct {
	Offset int64 `form:"offset" json:"offset"`
	Limit  int64 `form:"limit" json:"limit"`
}

type SyncMsgResponse struct {
	List []*Msg `json:"list"`
}

type MsgAckRequest struct {
	Id int64 `json:"id"`
}
