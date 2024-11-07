package zero

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
		atRegexp := regexp.MustCompile(`@([\d\S]*)`)
		for i := 0; i < len(ctx.Event.Message); i++ {
			if i > 0 && ctx.Event.Message[i-1].Type == "reply" && ctx.Event.Message[i].Type == "at" {
				// [reply][at]
				reply := ctx.GetMessage(ctx.Event.Message[i-1].Data["id"])
				if reply.MessageID.ID() != 0 && reply.Sender != nil && reply.Sender.ID != 0 && strconv.FormatInt(reply.Sender.ID, 10) == ctx.Event.Message[i].Data["qq"] {
					continue
				}
			}
			if ctx.Event.Message[i].Type == "text" && atRegexp.MatchString(ctx.Event.Message[i].Data["text"]) {
				// xxxx @11232123 xxxxx
				splited := atRegexp.Split(ctx.Event.Message[i].Data["text"], -1)
				ats := atRegexp.FindAllStringSubmatch(ctx.Event.Message[i].Data["text"], -1)
				var tmp = make([]message.Segment, 0, len(splited)+len(ats))
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
					if err != nil {
						// assume is user name
						list := ctx.GetThisGroupMemberList().Array()
						for _, member := range list {
							if member.Get("card").Str == ats[i][1] || member.Get("nickname").Str == ats[i][1] {
								uid = member.Get("user_id").Int()
								break
							}
						}
					}
					tmp = append(tmp, message.At(uid))
				}
				msgs = append(msgs, tmp...)
				continue
			}
			msgs = append(msgs, ctx.Event.Message[i])
		}
		return patternMatch(ctx, *p, msgs)
	}
}

type Pattern struct {
	cleanRedundantAt bool
	fuzzyAt          bool
	segments         []PatternSegment
}

// PatternOption pattern option
type PatternOption struct {
	cleanRedundantAt bool
	fuzzyAt          bool
}

// NewPattern new pattern
// defaults:
//
//	cleanRedundantAt: true
//	fuzzyAt: false
func NewPattern(cleanRedundantAt ...PatternOption) *Pattern {
	option := PatternOption{
		cleanRedundantAt: true,
		fuzzyAt:          false,
	}
	if len(cleanRedundantAt) > 0 {
		option = cleanRedundantAt[0]
	}
	pattern := Pattern{
		cleanRedundantAt: option.cleanRedundantAt,
		fuzzyAt:          option.fuzzyAt,
		segments:         make([]PatternSegment, 0, 4),
	}
	return &pattern
}

type PatternSegment struct {
	typ      string
	optional bool
	parse    Parser
	DebugStr func() string
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

func (p *Pattern) Add(typ string, optional bool, parse Parser, debug func() string) *Pattern {
	pattern := &PatternSegment{
		typ:      typ,
		optional: optional,
		parse:    parse,
		DebugStr: debug,
	}
	p.segments = append(p.segments, *pattern)
	return p
}

// Text use regex to search a 'text' segment
func (p *Pattern) Text(regex string) *Pattern {
	p.Add("text", false, NewTextParser(regex), func() string {
		return fmt.Sprintf("regex(%s)", regex)
	})
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
	p.Add("at", false, NewAtParser(id...), func() string {
		return fmt.Sprintf("at(%v)", id)
	})
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
	p.Add("image", false, NewImageParser(), func() string {
		return "image"
	})
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
	p.Add("reply", false, NewReplyParser(), func() string {
		return "reply"
	})
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
	p.Add("any", false, NewAnyParser(), func() string {
		return "any"
	})
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
