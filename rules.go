package zero

import (
	"hash/crc64"
	"regexp"
	"strconv"
	"strings"

	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

type Rule func(ctx *Ctx) bool

// RuleWrapper wraps a rule
func RuleWrapper(rule Rule) Handler {
	return func(ctx *Ctx) {
		if rule(ctx) {
			return
		}
		ctx.Abort()
	}
}

// Type check the ctx.Event's type
func Type(type_ string) Handler {
	t := strings.SplitN(type_, "/", 3)
	return func(ctx *Ctx) {
		if len(t) > 0 && t[0] != ctx.Event.PostType {
			ctx.Abort()
			return
		}
		if len(t) > 1 && t[1] != ctx.Event.DetailType {
			ctx.Abort()
			return
		}
		if len(t) > 2 && t[2] != ctx.Event.SubType {
			ctx.Abort()
		}
	}
}

// PrefixRule check if the message has the prefix and trim the prefix
//
// 检查消息前缀
func PrefixRule(prefixes ...string) Handler {
	return func(ctx *Ctx) {
		if len(ctx.Event.Message) == 0 || ctx.Event.Message[0].Type != "text" { // 确保无空指针
			ctx.Abort()
			return
		}
		first := ctx.Event.Message[0]
		firstMessage := first.Data["text"]
		for _, prefix := range prefixes {
			if strings.HasPrefix(firstMessage, prefix) {
				ctx.State["prefix"] = prefix
				arg := strings.TrimLeft(firstMessage[len(prefix):], " ")
				if len(ctx.Event.Message) > 1 {
					arg += ctx.Event.Message[1:].ExtractPlainText()
				}
				ctx.State["args"] = arg
				return
			}
		}
		ctx.Abort()
	}
}

// SuffixRule check if the message has the suffix and trim the suffix
//
// 检查消息后缀
func SuffixRule(suffixes ...string) Handler {
	return func(ctx *Ctx) {
		mLen := len(ctx.Event.Message)
		if mLen <= 0 { // 确保无空指针
			ctx.Abort()
			return
		}
		last := ctx.Event.Message[mLen-1]
		if last.Type != "text" {
			ctx.Abort()
			return
		}
		lastMessage := last.Data["text"]
		for _, suffix := range suffixes {
			if strings.HasSuffix(lastMessage, suffix) {
				ctx.State["suffix"] = suffix
				arg := strings.TrimRight(lastMessage[:len(lastMessage)-len(suffix)], " ")
				if mLen >= 2 {
					arg += ctx.Event.Message[:mLen].ExtractPlainText()
				}
				ctx.State["args"] = arg
				return
			}
		}
		ctx.Abort()
	}
}

// CommandRule check if the message is a command and trim the command name
func CommandRule(commands ...string) Handler {
	return func(ctx *Ctx) {
		if len(ctx.Event.Message) == 0 || ctx.Event.Message[0].Type != "text" {
			ctx.Abort()
			return
		}
		first := ctx.Event.Message[0]
		firstMessage := first.Data["text"]
		if !strings.HasPrefix(firstMessage, BotConfig.CommandPrefix) {
			ctx.Abort()
			return
		}
		cmdMessage := firstMessage[len(BotConfig.CommandPrefix):]
		for _, command := range commands {
			if strings.HasPrefix(cmdMessage, command) {
				ctx.State["command"] = command
				arg := strings.TrimLeft(cmdMessage[len(command):], " ")
				if len(ctx.Event.Message) > 1 {
					arg += ctx.Event.Message[1:].ExtractPlainText()
				}
				ctx.State["args"] = arg
				return
			}
		}
		ctx.Abort()
	}
}

// RegexRule check if the message can be matched by the regex pattern
func RegexRule(regexPattern string) Handler {
	regex := regexp.MustCompile(regexPattern)
	return func(ctx *Ctx) {
		msg := ctx.MessageString()
		if matched := regex.FindStringSubmatch(msg); matched != nil {
			ctx.State["regex_matched"] = matched
			return
		}
		ctx.Abort()
	}
}

// ReplyRule check if the message is replying some message
func ReplyRule(messageID int64) Handler {
	return func(ctx *Ctx) {
		if len(ctx.Event.Message) == 0 {
			ctx.Abort()
			return
		}
		if ctx.Event.Message[0].Type != "reply" {
			ctx.Abort()
			return
		}
		if id, err := strconv.ParseInt(ctx.Event.Message[0].Data["id"], 10, 64); err == nil {
			if id != messageID {
				ctx.Abort()
			}
			return
		}
		c := crc64.New(crc64.MakeTable(crc64.ISO))
		c.Write(helper.StringToBytes(ctx.Event.Message[0].Data["id"]))
		if int64(c.Sum64()) != messageID {
			ctx.Abort()
		}
	}
}

// KeywordRule check if the message has a keyword or keywords
func KeywordRule(src ...string) Handler {
	return func(ctx *Ctx) {
		msg := ctx.MessageString()
		for _, str := range src {
			if strings.Contains(msg, str) {
				ctx.State["keyword"] = str
				return
			}
		}
		ctx.Abort()
	}
}

// FullMatchRule check if src has the same copy of the message
func FullMatchRule(src ...string) Handler {
	return func(ctx *Ctx) {
		msg := ctx.MessageString()
		for _, str := range src {
			if str == msg {
				ctx.State["matched"] = msg
				return
			}
		}
		ctx.Abort()
	}
}

// OnlyToMe only triggered in conditions of @bot or begin with the nicknames
func OnlyToMe(ctx *Ctx) {
	if ctx.Event.IsToMe {
		return
	}
	ctx.Abort()
}

// CheckUser only triggered by specific person
func CheckUser(userId ...int64) Handler {
	return func(ctx *Ctx) {
		for _, uid := range userId {
			if ctx.Event.UserID == uid {
				return
			}
		}
		ctx.Abort()
	}
}

// OnlyPrivate requires that the ctx.Event is private message
func OnlyPrivate(ctx *Ctx) {
	if ctx.Event.PostType == "message" && ctx.Event.DetailType == "private" {
		return
	}
	ctx.Abort()
}

// OnlyPublic requires that the ctx.Event is public/group or public/guild message
func OnlyPublic(ctx *Ctx) {
	if ctx.Event.PostType == "message" && (ctx.Event.DetailType == "group" || ctx.Event.DetailType == "guild") {
		return
	}
	ctx.Abort()
}

// OnlyGroup requires that the ctx.Event is public/group message
func OnlyGroup(ctx *Ctx) {
	if ctx.Event.PostType == "message" && ctx.Event.DetailType == "group" {
		return
	}
	ctx.Abort()
}

// OnlyGuild requires that the ctx.Event is public/guild message
func OnlyGuild(ctx *Ctx) {
	if ctx.Event.PostType == "message" && ctx.Event.DetailType == "guild" {
		return
	}
	ctx.Abort()
}

// SuperUserPermission only triggered by the bot's owner
func SuperUserPermission(ctx *Ctx) {
	uid := strconv.FormatInt(ctx.Event.UserID, 10)
	for _, su := range BotConfig.SuperUsers {
		if su == uid {
			return
		}
	}
	ctx.Abort()
}

// AdminPermission only triggered by the group admins or higher permission
func AdminPermission(ctx *Ctx) {
	if ctx.Event.Sender.Role != "member" {
		return
	}
	SuperUserPermission(ctx)
}

// OwnerPermission only triggered by the group owner or higher permission
func OwnerPermission(ctx *Ctx) {
	if ctx.Event.Sender.Role == "owner" {
		return
	}
	SuperUserPermission(ctx)
}
