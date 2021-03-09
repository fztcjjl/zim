package router

import (
	"github.com/fztcjjl/zim/logic/handler/http/controller"
	"github.com/fztcjjl/zim/logic/handler/http/router/helper"
	"github.com/gin-gonic/gin"
)

func RegisterGroupRouter(g *gin.RouterGroup, c *controller.GroupController) {
	helper.POST(g, "/group/list", c.GetGroupList, "获取 App 中的所有群组")
	helper.POST(g, "/group/create", c.CreateGroup, "创建群组")
	helper.POST(g, "/group/get_group_info", c.GetGroupInfo, "获取群组信息")
	helper.POST(g, "/group/get_group_member_info", c.GetGroupMemberInfo, "获取群成员详细资料")
	helper.POST(g, "/group/modify_group_base_info", c.ModifyGroupBaseInfo, "修改群基础资料")
	helper.POST(g, "/group/add_group_member", c.AddGroupMember, "增加群成员")
	helper.POST(g, "/group/delete_group_member", c.DeleteGroupMember, "删除群成员")
	helper.POST(g, "/group/modify_group_member_info", c.ModifyGroupMemberInfo, "修改群成员资料")
	helper.POST(g, "/group/destroy_group", c.DestoryGroup, "解散群组")
	helper.POST(g, "/group/get_joined_group_list", c.GetJoinedGroupList, "获取用户所加入的群组")
	helper.POST(g, "/group/get_role_in_group", c.GetRoleInGroup, "查询用户在群组中的身份")
	helper.POST(g, "/group/forbid_send_msg", c.ForbidSendMsg, "批量禁言和取消禁言")
	helper.POST(g, "/group/get_group_shutted_uin", c.GetGroupShuttedUin, "获取被禁言群成员列表")
	helper.POST(g, "/group/send_group_msg", c.SendGroupMsg, "在群组中发送普通消息")
	helper.POST(g, "/group/send_group_system_notification", c.SendGroupSystemNotification, "在群组中发送系统通知")
	helper.POST(g, "/group/get_role_in_group", c.GetRoleInGroup, "查询用户在群组中的身份")
	helper.POST(g, "/group/group_msg_recall", c.GroupMsgRecall, "撤回群消息")
	helper.POST(g, "/group/change_group_owner", c.ChangeGroupOwner, "转让群主")
	helper.POST(g, "/group/group_msg_get", c.GroupMsgGet, "拉取群历史消息")
}
