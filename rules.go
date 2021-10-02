package zero

import (
	"regexp"
	"strconv"
	"strings"
)

// Type check the ctx.Event's type
func Type(type_ string) Rule {
	t := strings.SplitN(type_, "/", 3)
	return func(ctx *Ctx) bool {
		if len(t) > 0 && t[0] != ctx.Event.PostType {
			return false
		}
		if len(t) > 1 && t[1] != ctx.Event.DetailType {
			return false
		}
		if len(t) > 2 && t[2] != ctx.Event.SubType {
			return false
		}
		return true
	}
}

// PrefixRule check if the message has the prefix and trim the prefix
//
// 检查消息前缀
func PrefixRule(prefixes ...string) Rule {
	return func(ctx *Ctx) bool {
		if len(ctx.Event.Message) == 0 || ctx.Event.Message[0].Type != "text" { // 确保无空指针
			return false
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
				return true
			}
		}
		return false
	}
}

// SuffixRule check if the message has the suffix and trim the suffix
//
// 检查消息后缀
func SuffixRule(suffixes ...string) Rule {
	return func(ctx *Ctx) bool {
		mLen := len(ctx.Event.Message)
		if mLen <= 0 { // 确保无空指针
			return false
		}
		last := ctx.Event.Message[mLen-1]
		if last.Type != "text" {
			return false
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
				return true
			}
		}
		return false
	}
}

// CommandRule check if the message is a command and trim the command name
func CommandRule(commands ...string) Rule {
	return func(ctx *Ctx) bool {
		if len(ctx.Event.Message) == 0 || ctx.Event.Message[0].Type != "text" {
			return false
		}
		first := ctx.Event.Message[0]
		firstMessage := first.Data["text"]
		if !strings.HasPrefix(firstMessage, BotConfig.CommandPrefix) {
			return false
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
				return true
			}
		}
		return false
	}
}

// RegexRule check if the message can be matched by the regex pattern
func RegexRule(regexPattern string) Rule {
	regex := regexp.MustCompile(regexPattern)
	return func(ctx *Ctx) bool {
		msg := ctx.MessageString()
		if regex.MatchString(msg) {
			ctx.State["regex_matched"] = regex.FindStringSubmatch(msg)
			return true
		}
		return false
	}
}

// ReplyRule check if the message is replying some message
func ReplyRule(messageID int64) Rule {
	mid := strconv.FormatInt(messageID, 10)
	return func(ctx *Ctx) bool {
		if len(ctx.Event.Message) == 0 {
			return false
		}
		if ctx.Event.Message[0].Type != "reply" {
			return false
		}
		return ctx.Event.Message[0].Data["id"] == mid
	}
}

// KeywordRule check if the message has a keyword or keywords
func KeywordRule(src ...string) Rule {
	return func(ctx *Ctx) bool {
		msg := ctx.MessageString()
		for _, str := range src {
			if strings.Contains(msg, str) {
				ctx.State["keyword"] = str
				return true
			}
		}
		return false
	}
}

// FullMatchRule check if src has the same copy of the message
func FullMatchRule(src ...string) Rule {
	return func(ctx *Ctx) bool {
		msg := ctx.MessageString()
		for _, str := range src {
			if str == msg {
				ctx.State["matched"] = msg
				return true
			}
		}
		return false
	}
}

// OnlyToMe only triggered in conditions of @bot or begin with the nicknames
func OnlyToMe(ctx *Ctx) bool {
	return ctx.Event.IsToMe
}

// CheckUser only triggered by specific person
func CheckUser(userId ...int64) Rule {
	return func(ctx *Ctx) bool {
		for _, uid := range userId {
			if ctx.Event.UserID == uid {
				return true
			}
		}
		return false
	}
}

// OnlyPrivate requires that the ctx.Event is private message
func OnlyPrivate(ctx *Ctx) bool {
	return ctx.Event.PostType == "message" && ctx.Event.DetailType == "private"
}

// OnlyGroup requires that the ctx.Event is public/group message
func OnlyGroup(ctx *Ctx) bool {
	return ctx.Event.PostType == "message" && ctx.Event.DetailType == "group"
}

// SuperUserPermission only triggered by the bot's owner
func SuperUserPermission(ctx *Ctx) bool {
	uid := strconv.FormatInt(ctx.Event.UserID, 10)
	for _, su := range BotConfig.SuperUsers {
		if su == uid {
			return true
		}
	}
	return false
}

// AdminPermission only triggered by the group admins or higher permission
func AdminPermission(ctx *Ctx) bool {
	return SuperUserPermission(ctx) || ctx.Event.Sender.Role != "member"
}

// OwnerPermission only triggered by the group owner or higher permission
func OwnerPermission(ctx *Ctx) bool {
	return SuperUserPermission(ctx) ||
		(ctx.Event.Sender.Role != "member" && ctx.Event.Sender.Role != "admin")
}
