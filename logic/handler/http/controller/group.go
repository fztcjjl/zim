package controller

import "github.com/fztcjjl/zim/pkg/gin/context"

type GroupController struct {
}

// 获取 App 中的所有群组
func (t *GroupController) GetGroupList(ctx *context.Context) {
	//ctx.GetHeader()
}

// 创建群组
func (t *GroupController) CreateGroup(ctx *context.Context) {

}

// 获取群详细资料
func (t *GroupController) GetGroupInfo(ctx *context.Context) {

}

// 解散群组
func (t *GroupController) DestoryGroup(ctx *context.Context) {

}

// 获取群成员详细资料
func (t *GroupController) GetGroupMemberInfo(ctx *context.Context) {

}

// 修改群基础资料
func (t *GroupController) ModifyGroupBaseInfo(ctx *context.Context) {

}

// 增加群成员
func (t *GroupController) AddGroupMember(ctx *context.Context) {

}

// 删除群成员
func (t *GroupController) DeleteGroupMember(ctx *context.Context) {

}

// 修改群成员资料
func (t *GroupController) ModifyGroupMemberInfo(ctx *context.Context) {

}

// 获取用户所加入的群组
func (t *GroupController) GetJoinedGroupList(ctx *context.Context) {

}

// 查询用户在群组中的身份
func (t *GroupController) GetRoleInGroup(ctx *context.Context) {

}

// 批量禁言和取消禁言
func (t *GroupController) ForbidSendMsg(ctx *context.Context) {

}

// 获取被禁言群成员列表
func (t *GroupController) GetGroupShuttedUin(ctx *context.Context) {

}

// 在群组中发送普通消息
func (t *GroupController) SendGroupMsg(ctx *context.Context) {

}

// 在群组中发送系统通知
func (t *GroupController) SendGroupSystemNotification(ctx *context.Context) {

}

// 撤回群消息
func (t *GroupController) GroupMsgRecall(ctx *context.Context) {

}

// 转让群主
func (t *GroupController) ChangeGroupOwner(ctx *context.Context) {

}

// 拉取群历史消息
func (t *GroupController) GroupMsgGet(ctx *context.Context) {

}
