# ZeroBot

[![Go Report Card](https://goreportcard.com/badge/github.com/wdvxdr1123/ZeroBot)](https://goreportcard.com/report/github.com/wdvxdr1123/ZeroBot)
![golangci-lint](https://github.com/wdvxdr1123/ZeroBot/workflows/golang-ci/badge.svg)
![Badge](https://img.shields.io/badge/OneBot-v11-black)
![Badge](https://img.shields.io/badge/gocqhttp-v1.0.0-black)
[![License](https://img.shields.io/github/license/wdvxdr1123/ZeroBot.svg?style=flat-square&logo=gnu)](https://raw.githubusercontent.com/wdvxdr1123/ZeroBot/main/LICENSE)
[![qq group](https://img.shields.io/badge/group-892659456-red?style=flat-square&logo=tencent-qq)](https://jq.qq.com/?_wv=1027&k=E6Zov6Fi)



## ⚡️ 快速使用

```go
package main

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	zero.OnCommand("hello").
            Handle(func(ctx *zero.Ctx) {
                ctx.Send("world")
            })

	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{123456},
		Driver: []zero.Driver{
			// 正向 WS
			driver.NewWebSocketClient("ws://127.0.0.1:6700", ""),
			// 反向 WS
			driver.NewWebSocketServer(16, "ws://127.0.0.1:6701", ""),
			// HTTP
			driver.NewHTTPClient("http://127.0.0.1:6701", "", "http://127.0.0.1:6700", ""),
		},
	}, nil)
}
```

## 📖 文档

[**--> 点击这里查看文档 <--**](https://wdvxdr1123.github.io/ZeroBot/)

## 🎯 特性

- 通过 `init` 函数实现插件式
- 底层与 Onebot 通信驱动可换，目前支持HTTP、正向/反向WS，且支持基于 `unix socket` 的通信（使用 `ws+unix://`）
- 通过添加多个 driver 实现多Q机器人支持
- 完整的 OneBot 11 标准 API + NapCat / NapNeko 扩展 API 支持

## 📡 API 列表

所有 API 方法均定义在 `api.go` 中，以下按功能分组列出。带 `ThisGroup` / `This` 前缀的方法是对应方法的便捷版本，自动使用当前事件的群号/消息ID。

---

### OneBot 11 标准 API

#### 消息

| 方法 | 说明 | Action |
|------|------|--------|
| `SendGroupMessage` | 发送群消息 | `send_group_msg` |
| `SendPrivateMessage` | 发送私聊消息 | `send_private_msg` |
| `DeleteMessage` | 撤回消息 | `delete_msg` |
| `GetMessage` | 获取消息 | `get_msg` |
| `GetForwardMessage` | 获取合并转发消息 | `get_forward_msg` |
| `SendLike` | 发送好友赞 | `send_like` |

#### 群操作

| 方法 | 说明 | Action |
|------|------|--------|
| `SetGroupKick` / `SetThisGroupKick` | 群组踢人 | `set_group_kick` |
| `SetGroupBan` / `SetThisGroupBan` | 群组单人禁言 | `set_group_ban` |
| `SetGroupWholeBan` / `SetThisGroupWholeBan` | 群组全员禁言 | `set_group_whole_ban` |
| `SetGroupAdmin` / `SetThisGroupAdmin` | 群组设置管理员 | `set_group_admin` |
| `SetGroupAnonymous` / `SetThisGroupAnonymous` | 群组匿名 | `set_group_anonymous` |
| `SetGroupCard` / `SetThisGroupCard` | 设置群名片 | `set_group_card` |
| `SetGroupName` / `SetThisGroupName` | 设置群名 | `set_group_name` |
| `SetGroupLeave` / `SetThisGroupLeave` | 退出群组 | `set_group_leave` |
| `SetGroupSpecialTitle` / `SetThisGroupSpecialTitle` | 设置群组专属头衔 | `set_group_special_title` |

#### 请求处理

| 方法 | 说明 | Action |
|------|------|--------|
| `SetFriendAddRequest` | 处理加好友请求 | `set_friend_add_request` |
| `SetGroupAddRequest` | 处理加群请求/邀请 | `set_group_add_request` |

#### 账号信息

| 方法 | 说明 | Action |
|------|------|--------|
| `GetLoginInfo` | 获取登录号信息 | `get_login_info` |
| `GetStrangerInfo` | 获取陌生人信息 | `get_stranger_info` |
| `GetFriendList` | 获取好友列表 | `get_friend_list` |
| `GetGroupInfo` / `GetThisGroupInfo` | 获取群信息 | `get_group_info` |
| `GetGroupList` | 获取群列表 | `get_group_list` |
| `GetGroupMemberInfo` / `GetThisGroupMemberInfo` | 获取群成员信息 | `get_group_member_info` |
| `GetGroupMemberList` / `GetThisGroupMemberList` | 获取群成员列表 | `get_group_member_list` |
| `GetGroupMemberListNoCache` / `GetThisGroupMemberListNoCache` | 无缓存获取群成员列表 | `get_group_member_list` |
| `GetGroupHonorInfo` / `GetThisGroupHonorInfo` | 获取群荣誉信息 | `get_group_honor_info` |
| `GetRecord` | 获取语音 | `get_record` |
| `GetImage` | 获取图片 | `get_image` |
| `GetVersionInfo` | 获取版本信息 | `get_version_info` |

---

### go-cqhttp 扩展 API

#### 群操作扩展

| 方法 | 说明 | Action |
|------|------|--------|
| `SetGroupPortrait` / `SetThisGroupPortrait` | 设置群头像 | `set_group_portrait` |
| `GetGroupSystemMessage` | 获取群系统消息 | `get_group_system_msg` |
| `GetGroupAtAllRemain` / `GetThisGroupAtAllRemain` | 获取群@全体成员剩余次数 | `get_group_at_all_remain` |
| `GetGroupMessageHistory` / `GetThisGroupMessageHistory` | 获取群消息历史记录 | `get_group_msg_history` |
| `GetLatestGroupMessageHistory` / `GetLatestThisGroupMessageHistory` | 获取最新群消息历史记录 | `get_group_msg_history` |
| `GetGroupEssenceMessageList` / `GetThisGroupEssenceMessageList` | 获取群精华消息列表 | `get_essence_msg_list` |
| `SetGroupEssenceMessage` | 设置群精华消息 | `set_essence_msg` |
| `DeleteGroupEssenceMessage` | 移出群精华消息 | `delete_essence_msg` |

#### 消息扩展

| 方法 | 说明 | Action |
|------|------|--------|
| `SendGroupForwardMessage` | 发送合并转发(群) | `send_group_forward_msg` |
| `SendPrivateForwardMessage` | 发送合并转发(私聊) | `send_private_forward_msg` |
| `ForwardFriendSingleMessage` | 转发单条消息到好友 | `forward_friend_single_msg` |
| `ForwardGroupSingleMessage` | 转发单条消息到群 | `forward_group_single_msg` |
| `MarkMessageAsRead` / `MarkThisMessageAsRead` | 标记消息已读 | `mark_msg_as_read` |
| `GetOnlineClients` | 获取当前账号在线客户端列表 | `get_online_clients` |

#### 文件操作

| 方法 | 说明 | Action |
|------|------|--------|
| `GetGroupFilesystemInfo` / `GetThisGroupFilesystemInfo` | 获取群文件系统信息 | `get_group_file_system_info` |
| `GetGroupRootFiles` / `GetThisGroupRootFiles` | 获取群根目录文件列表 | `get_group_root_files` |
| `GetGroupFilesByFolder` / `GetThisGroupFilesByFolder` | 获取群子目录文件列表 | `get_group_files_by_folder` |
| `GetGroupFileURL` / `GetThisGroupFileURL` | 获取群文件资源链接 | `get_group_file_url` |
| `UploadGroupFile` / `UploadThisGroupFile` | 上传群文件 | `upload_group_file` |

#### 其他

| 方法 | 说明 | Action |
|------|------|--------|
| `OCRImage` | 图片 OCR | `ocr_image` |
| `GetWordSlices` | 获取中文分词 | `.get_word_slices` |
| `SendGuildChannelMessage` | 发送频道消息 | `send_guild_channel_msg` |
| `NickName` | 从 args/at 获取昵称 | — |
| `CardOrNickName` | 从 uid 获取群名片或昵称 | — |

---

### NapNeko / LLOneBot 扩展 API

| 方法 | 说明 | Action |
|------|------|--------|
| `SetMyAvatar` | 设置我的头像 | `set_qq_avatar` |
| `GetFile` | 下载收到的文件 | `get_file` |
| `SetMessageEmojiLike` | 发送表情回应 | `set_msg_emoji_like` |
| `SetGroupSign` | 群签到 | `set_group_sign` |
| `GroupPoke` | 群聊戳一戳 | `group_poke` |
| `FriendPoke` | 私聊戳一戳 | `friend_poke` |
| `SendPoke` | 群聊/私聊戳一戳 | `send_poke` |
| `ArkSharePeer` | 获取推荐好友/群聊卡片 | `ArkSharePeer` |
| `ArkShareGroup` | 获取推荐群聊卡片 | `ArkShareGroup` |
| `GetRobotUinRange` | 获取机器人账号范围 | `get_robot_uin_range` |
| `SetOnlineStatus` | 设置在线状态 | `set_online_status` |
| `GetFriendsWithCategory` | 获取分类的好友列表 | `get_friends_with_category` |
| `TranslateEn2Zh` | 英译中 | `translate_en2zh` |
| `SendForwardMessage` | 发送合并转发 | `send_forward_msg` |
| `MarkPrivateMessageAsRead` | 设置私聊已读 | `mark_private_msg_as_read` |
| `MarkGroupMessageAsRead` | 设置群聊已读 | `mark_group_msg_as_read` |
| `GetFriendMessageHistory` | 获取私聊历史记录 | `get_friend_msg_history` |
| `CreateCollection` | 创建收藏 | `create_collection` |
| `GetCollectionList` | 获取收藏 | `get_collection_list` |
| `SetSelfLongNick` | 设置签名 | `set_self_longnick` |
| `GetRecentContact` | 获取最近联系人 | `get_recent_contact` |
| `MarkAllAsRead` | 标记所有已读 | `_mark_all_as_read` |
| `GetProfileLike` | 获取自身点赞列表 | `get_profile_like` |
| `FetchCustomFace` | 获取自定义表情 | `fetch_custom_face` |
| `GetAIRecord` | AI 文字转语音 | `get_ai_record` |
| `GetAICharacters` | 获取 AI 语音角色列表 | `get_ai_characters` |
| `SendGroupAIRecord` | 群聊发送 AI 语音 | `send_group_ai_record` |

---

### NapCat 扩展 API

基于 [NapCat API 文档](https://napcat.apifox.cn) 补充的扩展接口：

#### 文件操作

| 方法 | 说明 | Action |
|------|------|--------|
| `UploadPrivateFile` | 上传私聊文件 | `upload_private_file` |
| `DeleteGroupFile` / `DeleteThisGroupFile` | 删除群文件 | `delete_group_file` |
| `CreateGroupFileFolder` / `CreateThisGroupFileFolder` | 创建群文件目录 | `create_group_file_folder` |
| `DeleteGroupFileFolder` / `DeleteThisGroupFileFolder` | 删除群文件目录 | `delete_group_folder` |
| `DownloadFile` | 下载文件到本地 | `download_file` |
| `GetPrivateFileURL` | 获取私聊文件下载链接 | `get_private_file_url` |
| `MoveGroupFile` | 移动群文件 | `move_group_file` |
| `RenameGroupFile` | 重命名群文件 | `rename_group_file` |
| `TransGroupFile` | 传输群文件 | `trans_group_file` |

#### 群公告

| 方法 | 说明 | Action |
|------|------|--------|
| `SendGroupNotice` / `SendThisGroupNotice` | 发送群公告 | `_send_group_notice` |
| `GetGroupNotice` / `GetThisGroupNotice` | 获取群公告列表 | `_get_group_notice` |
| `DeleteGroupNotice` / `DeleteThisGroupNotice` | 删除群公告 | `_del_group_notice` |

#### 好友管理

| 方法 | 说明 | Action |
|------|------|--------|
| `DeleteFriend` | 删除好友 | `delete_friend` |
| `SetFriendRemark` | 设置好友备注 | `set_friend_remark` |
| `GetUnidirectionalFriendList` | 获取单向好友列表 | `get_unidirectional_friend_list` |
| `GetDoubtFriendsAddRequest` | 获取可疑好友申请 | `get_doubt_friends_add_request` |
| `SetDoubtFriendsAddRequest` | 处理可疑好友申请 | `set_doubt_friends_add_request` |

#### 用户资料

| 方法 | 说明 | Action |
|------|------|--------|
| `SetQQProfile` | 设置QQ资料 | `set_qq_profile` |
| `SetInputStatus` | 设置输入状态 | `set_input_status` |
| `GetUserStatus` | 获取用户在线状态 | `nc_get_user_status` |
| `SetCustomOnlineStatus` | 设置自定义在线状态 | `set_custom_online_status` |

#### 群管理扩展

| 方法 | 说明 | Action |
|------|------|--------|
| `GetGroupShutList` / `GetThisGroupShutList` | 获取群禁言列表 | `get_group_shut_list` |
| `GetGroupInfoEx` / `GetThisGroupInfoEx` | 获取群详细信息 | `get_group_info_ex` |
| `SetGroupRemark` / `SetThisGroupRemark` | 设置群备注 | `set_group_remark` |
| `GetGroupIgnoredNotifies` / `GetThisGroupIgnoredNotifies` | 获取群忽略通知 | `get_group_ignored_notifies` |
| `SetGroupAddOption` | 设置群加群选项 | `set_group_add_option` |
| `SetGroupSearchOption` | 设置群搜索选项 | `set_group_search_option` |
| `GroupKickBatch` | 批量踢出群成员 | `set_group_kick` |
| `SetGroupTodo` / `SetThisGroupTodo` | 设置群待办 | `set_group_todo` |

#### 群相册

| 方法 | 说明 | Action |
|------|------|--------|
| `GetGroupAlbumList` | 获取群相册列表 | `get_group_album_list` |
| `GetGroupAlbumMediaList` | 获取群相册媒体列表 | `get_group_album_media_list` |
| `UploadGroupAlbum` | 上传图片到群相册 | `upload_group_album` |
| `DeleteGroupAlbumMedia` | 删除群相册媒体 | `delete_group_album_media` |
| `LikeGroupAlbumMedia` | 点赞群相册媒体 | `like_group_album_media` |

#### 消息扩展

| 方法 | 说明 | Action |
|------|------|--------|
| `SendGroupMusic` | 发送群聊音乐卡片 | `send_group_msg` |
| `SendGroupCustomMusic` | 发送群聊自定义音乐卡片 | `send_group_msg` |
| `GetEmojiLikeList` | 获取消息表情点赞列表 | `get_msg_emoji_like_list` |
| `GetMiniAppArk` | 获取小程序 Ark | `get_mini_app_ark` |

#### 系统 / 安全

| 方法 | 说明 | Action |
|------|------|--------|
| `CheckURLSafely` | 检查URL安全性 | `check_url_safely` |
| `CanSendImage` | 检查是否可以发送图片 | `can_send_image` |
| `CanSendRecord` | 检查是否可以发送语音 | `can_send_record` |
| `GetCSRFToken` | 获取 CSRF Token | `get_csrf_token` |
| `GetCredentials` | 获取登录凭证 | `get_credentials` |
| `GetCookies` | 获取 Cookies | `get_cookies` |
| `GetClientKey` | 获取 ClientKey | `get_clientkey` |
| `GetStatus` | 获取运行状态 | `get_status` |
| `CleanCache` | 清理缓存 | `clean_cache` |
| `Restart` | 重启服务 | `set_restart` |
| `GetPacketStatus` | 获取Packet状态 | `get_packet_status` |
| `Logout` | 退出登录 | `nc_logout` |

#### 频道 / RKey / 其他

| 方法 | 说明 | Action |
|------|------|--------|
| `GetGuildList` | 获取频道列表 | `get_guild_list` |
| `GetGuildServiceProfile` | 获取频道个人信息 | `get_guild_service_profile` |
| `GetRKey` | 获取 RKey | `get_rkey` |
| `NcGetRKey` | 获取扩展 RKey | `nc_get_rkey` |
| `GetRKeyServer` | 获取 RKey 服务器 | `get_rkey_server` |
| `ClickInlineKeyboardButton` | 点击内联键盘按钮 | `click_inline_keyboard_button` |

## 关联项目

- [ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin): 基于 ZeroBot 的 OneBot 插件合集

## 特别感谢

- [nonebot/nonebot2](https://github.com/nonebot/nonebot2): 跨平台 Python 异步聊天机器人框架

- [catsworld/qq-bot-api](https://github.com/catsworld/qq-bot-api): Golang bindings for the Coolq HTTP API Plugin


同时感谢以下开发者对 ZeroBot 作出的贡献：

<a href="https://github.com/wdvxdr1123/ZeroBot/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=wdvxdr1123/ZeroBot&max=1000" />
</a>
