package api

type TextMsg struct {
	Text string `json:"text"`
}

//int32 conv_type = 1;
//string from = 2;
//string to = 3;
//string target = 4;
//string content = 5;
//int64 send_time = 6;
//string extra = 7;

type SendRequest struct {
	ConvType int    `form:"conv_type" json:"conv_type" binding:"required"`
	MsgType  int    `form:"msg_type" json:"msg_type" binding:"required"`
	From     string `form:"from" json:"from" binding:"required"`
	To       string `form:"to" json:"to" binding:"required"`
	Body     string `form:"body" json:"body" binding:"required"`
	Extra    string `form:"extra" json:"extra"`
}

type Msg struct {
	Id       int64  `json:"id"`
	ConvType int    `json:"conv_type"`
	Type     int    `json:"type"`
	Body     string `json:"body"`
	Extra    string `json:"extra"`
	From     string `json:"from"`
	To       string `json:"to"`
	SendTime int64  `json:"send_time"`
	ReadTime int64  `json:"read_time"`
}

type SyncMsgRequest struct {
	Limit  int   `form:"limit" json:"limit"`
	Offset int64 `form:"offset" json:"offset"`
}

type SyncMsgResponse struct {
	List []*Msg `json:"list"`
}

type MsgAckRequest struct {
	Id int64 `json:"id"`
}
