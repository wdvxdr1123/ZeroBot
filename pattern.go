package zero

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	KeyPattern = "pattern_matched"
)

// AsRule build PatternRule
func (p *Pattern) AsRule() Rule {
	return func(ctx *Ctx) bool {
		if len(ctx.Event.Message) == 0 {
			return false
		}
		if !p.cleanRedundantAt && !p.fuzzyAt {
			return patternMatch(ctx, *p, ctx.Event.Message)
		}

		// copy messages
		msgs := make([]message.Segment, 0, len(ctx.Event.Message))
		for i := 0; i < len(ctx.Event.Message); i++ {
			if i > 0 && ctx.Event.Message[i-1].Type == "reply" && ctx.Event.Message[i].Type == "at" {
				// [reply][at]
				// use internal API to avoid recording wrong triggered message
				reply := Message{}
				msgID := ctx.Event.Message[i-1].Data["id"]
				rsp, err := ctx.caller.(*messageLogger).caller.CallAPI(APIRequest{
					Action: "get_msg", Params: Params{
						"message_id": msgID,
					},
				})
				if err != nil {
					log.Debugf("[PatternRule] failed to get_msg, message_id %s", msgID)
					continue
				}
				reply = Message{
					Elements:    message.ParseMessage(helper.StringToBytes(rsp.Data.Get("message").Raw)),
					MessageID:   message.NewMessageIDFromInteger(rsp.Data.Get("message_id").Int()),
					MessageType: rsp.Data.Get("message_type").String(),
					Sender:      &User{},
				}
				_ = json.Unmarshal(helper.StringToBytes(rsp.Data.Get("sender").Raw), reply.Sender)
				if reply.MessageID.ID() != 0 && reply.Sender != nil && reply.Sender.ID != 0 && strconv.FormatInt(reply.Sender.ID, 10) == ctx.Event.Message[i].Data["qq"] {
					continue
				}
			}
			if ctx.Event.Message[i].Type == "text" && atRegexp.MatchString(ctx.Event.Message[i].Data["text"]) {
				// xxxx @11232123 xxxxx
				msgs = append(msgs, ctx.splitAtInText(i)...)
				continue
			}
			msgs = append(msgs, ctx.Event.Message[i])
		}
		return patternMatch(ctx, *p, msgs)
	}
}

var atRegexp = regexp.MustCompile(`@([\d\S]*)`)

func (ctx *Ctx) splitAtInText(index int) []message.Segment {
	msg := ctx.Event.Message[index].String()
	splited := atRegexp.Split(msg, -1)
	ats := atRegexp.FindAllStringSubmatch(msg, -1)
	var tmp = make([]message.Segment, 0, len(splited)+len(ats))
	var list []gjson.Result
	for i, s := range splited {
		if strings.TrimSpace(s) == "" {
			continue
		}
		tmp = append(tmp, message.Text(s))
		// append at
		if i > len(ats)-1 {
			continue
		}
		uid, err := strconv.ParseInt(ats[i][1], 10, 64)
		// TODO numeric username
		if err != nil {
			// assume is username
			if list == nil {
				list = ctx.GetThisGroupMemberList().Array()
			}
			for _, member := range list {
				if member.Get("card").Str != ats[i][1] && member.Get("nickname").Str != ats[i][1] {
					continue
				}
				uid = member.Get("user_id").Int()
			}
		}
		tmp = append(tmp, message.At(uid))
	}
	return tmp
}

type Pattern struct {
	cleanRedundantAt bool
	fuzzyAt          bool
	segments         []PatternSegment
}

// PatternOption pattern option
type PatternOption struct {
	CleanRedundantAt bool
	FuzzyAt          bool
}

// NewPattern new pattern
// defaults:
//
//	CleanRedundantAt: true
//	FuzzyAt: false
func NewPattern(option *PatternOption) *Pattern {
	if option == nil {
		option = &PatternOption{
			CleanRedundantAt: true,
			FuzzyAt:          false,
		}
	}
	pattern := Pattern{
		cleanRedundantAt: option.CleanRedundantAt,
		fuzzyAt:          option.FuzzyAt,
		segments:         make([]PatternSegment, 0, 4),
	}
	return &pattern
}

type PatternSegment struct {
	typ      string
	optional bool
	parse    Parser
}

type Parser func(msg *message.Segment) PatternParsed

// SetOptional set previous segment is optional, is v is empty, optional will be true
// if Pattern is empty, panic
func (p *Pattern) SetOptional(v ...bool) *Pattern {
	if len(p.segments) == 0 {
		panic("pattern is empty")
	}
	if len(v) == 1 {
		p.segments[len(p.segments)-1].optional = v[0]
	} else {
		p.segments[len(p.segments)-1].optional = true
	}
	return p
}

// PatternParsed PatternRule parse result
type PatternParsed struct {
	value any
	msg   *message.Segment
}

// Text 获取正则表达式匹配到的文本数组
func (p PatternParsed) Text() []string {
	if p.value == nil {
		return nil
	}
	return p.value.([]string)
}

// At 获取被@者ID
func (p PatternParsed) At() string {
	if p.value == nil {
		return ""
	}
	return p.value.(string)
}

// Image 获取图片URL
func (p PatternParsed) Image() string {
	if p.value == nil {
		return ""
	}
	return p.value.(string)
}

// Reply 获取被回复的消息ID
func (p PatternParsed) Reply() string {
	if p.value == nil {
		return ""
	}
	return p.value.(string)
}

// Raw 获取原始消息
func (p PatternParsed) Raw() *message.Segment {
	return p.msg
}

func (p *Pattern) Add(typ string, optional bool, parse Parser) *Pattern {
	pattern := &PatternSegment{
		typ:      typ,
		optional: optional,
		parse:    parse,
	}
	p.segments = append(p.segments, *pattern)
	return p
}

// Text use regex to search a 'text' segment
func (p *Pattern) Text(regex string) *Pattern {
	p.Add("text", false, NewTextParser(regex))
	return p
}

func NewTextParser(regex string) Parser {
	re := regexp.MustCompile(regex)
	return func(msg *message.Segment) PatternParsed {
		s := msg.Data["text"]
		s = strings.Trim(s, " \n\r\t")
		matchString := re.MatchString(s)
		if matchString {
			return PatternParsed{
				value: re.FindStringSubmatch(s),
				msg:   msg,
			}
		}

		return PatternParsed{}
	}
}

// At use regex to match an 'at' segment, if id is not empty, only match specific target
func (p *Pattern) At(id ...message.ID) *Pattern {
	if len(id) > 1 {
		panic("at pattern only support one id")
	}
	p.Add("at", false, NewAtParser(id...))
	return p
}

func NewAtParser(id ...message.ID) Parser {
	return func(msg *message.Segment) PatternParsed {
		if len(id) == 0 || len(id) == 1 && id[0].String() == msg.Data["qq"] {
			return PatternParsed{
				value: msg.Data["qq"],
				msg:   msg,
			}
		}
		return PatternParsed{}
	}
}

// Image use regex to match an 'at' segment, if id is not empty, only match specific target
func (p *Pattern) Image() *Pattern {
	p.Add("image", false, NewImageParser())
	return p
}

func NewImageParser() Parser {
	return func(msg *message.Segment) PatternParsed {
		return PatternParsed{
			value: msg.Data["file"],
			msg:   msg,
		}
	}
}

// Reply type zero.PatternReplyMatched
func (p *Pattern) Reply() *Pattern {
	p.Add("reply", false, NewReplyParser())
	return p
}

func NewReplyParser() Parser {
	return func(msg *message.Segment) PatternParsed {
		return PatternParsed{
			value: msg.Data["id"],
			msg:   msg,
		}
	}
}

// Any match any segment
func (p *Pattern) Any() *Pattern {
	p.Add("any", false, NewAnyParser())
	return p
}

func NewAnyParser() Parser {
	return func(msg *message.Segment) PatternParsed {
		parsed := PatternParsed{
			value: nil,
			msg:   msg,
		}
		switch {
		case msg.Data["text"] != "":
			parsed.value = msg.Data["text"]
		case msg.Data["qq"] != "":
			parsed.value = msg.Data["qq"]
		case msg.Data["file"] != "":
			parsed.value = msg.Data["file"]
		case msg.Data["id"] != "":
			parsed.value = msg.Data["id"]
		default:
			parsed.value = msg.Data
		}
		return parsed
	}
}

func (s *PatternSegment) matchType(msg message.Segment) bool {
	return s.typ == msg.Type || s.typ == "any"
}
func mustMatchAllPatterns(pattern Pattern) bool {
	for _, p := range pattern.segments {
		if p.optional {
			return false
		}
	}
	return true
}
func patternMatch(ctx *Ctx, pattern Pattern, msgs []message.Segment) bool {
	if mustMatchAllPatterns(pattern) && len(pattern.segments) != len(msgs) {
		return false
	}
	patternState := make([]PatternParsed, len(pattern.segments))

	j := 0
	for i := range pattern.segments {
		if j < len(msgs) && pattern.segments[i].matchType(msgs[j]) {
			patternState[i] = pattern.segments[i].parse(&msgs[j])
		}
		if patternState[i].value == nil {
			if pattern.segments[i].optional {
				continue
			}
			return false
		}
		j++
	}
	ctx.State[KeyPattern] = patternState
	return true
}
