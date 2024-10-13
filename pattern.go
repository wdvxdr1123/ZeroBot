package zero

import (
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
		// copy messages
		msgs := make([]message.Segment, 0, len(ctx.Event.Message))
		msgs = append(msgs, ctx.Event.Message[0])
		shouldClean := false
		for _, segment := range *p {
			if segment.cleanRedundantAt {
				shouldClean = true
				break
			}
		}
		for i := 1; i < len(ctx.Event.Message); i++ {
			if shouldClean && ctx.Event.Message[i-1].Type == "reply" && ctx.Event.Message[i].Type == "at" {
				// [reply][at]
				reply := ctx.GetMessage(ctx.Event.Message[i-1].Data["id"])
				if reply.MessageID.ID() != 0 && reply.Sender != nil && reply.Sender.ID != 0 && strconv.FormatInt(reply.Sender.ID, 10) == ctx.Event.Message[i].Data["qq"] {
					continue
				}
			}
			msgs = append(msgs, ctx.Event.Message[i])
		}
		return patternMatch(ctx, *p, msgs)
	}
}

type Pattern []PatternSegment

func NewPattern() *Pattern {
	pattern := make(Pattern, 0, 4)
	return &pattern
}

type PatternSegment struct {
	typ              string
	optional         bool
	parse            func(msg *message.Segment) *PatternParsed
	cleanRedundantAt bool // only for Reply
}

// SetOptional set previous segment is optional, is v is empty, optional will be true
// if Pattern is empty, panic
func (p *Pattern) SetOptional(v ...bool) *Pattern {
	if len(*p) == 0 {
		panic("pattern is empty")
	}
	if len(v) == 1 {
		(*p)[len(*p)-1].optional = v[0]
	} else {
		(*p)[len(*p)-1].optional = true
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

func NewPatternSegment(typ string, optional bool, parse func(msg *message.Segment) *PatternParsed, cleanRedundantAt ...bool) *PatternSegment {
	clean := false
	if len(cleanRedundantAt) > 0 {
		clean = cleanRedundantAt[0]
	}
	return &PatternSegment{
		typ:              typ,
		optional:         optional,
		parse:            parse,
		cleanRedundantAt: clean,
	}
}

// Text use regex to search a 'text' segment
func (p *Pattern) Text(regex string) *Pattern {
	re := regexp.MustCompile(regex)
	pattern := NewPatternSegment(
		"text", false, func(msg *message.Segment) *PatternParsed {
			s := msg.Data["text"]
			s = strings.Trim(s, " \n\r\t")
			matchString := re.MatchString(s)
			if matchString {
				return &PatternParsed{
					value: re.FindStringSubmatch(s),
					msg:   msg,
				}
			}

			return &PatternParsed{
				value: nil,
				msg:   nil,
			}
		},
	)
	*p = append(*p, *pattern)
	return p
}

// At use regex to match an 'at' segment, if id is not empty, only match specific target
func (p *Pattern) At(id ...message.ID) *Pattern {
	if len(id) > 1 {
		panic("at pattern only support one id")
	}
	pattern := NewPatternSegment(
		"at", false, func(msg *message.Segment) *PatternParsed {
			if len(id) == 0 || len(id) == 1 && id[0].String() == msg.Data["qq"] {
				return &PatternParsed{
					value: msg.Data["qq"],
					msg:   msg,
				}
			}

			return &PatternParsed{
				value: nil,
				msg:   nil,
			}
		},
	)
	*p = append(*p, *pattern)
	return p
}

// Image use regex to match an 'at' segment, if id is not empty, only match specific target
func (p *Pattern) Image() *Pattern {
	pattern := NewPatternSegment(
		"image", false, func(msg *message.Segment) *PatternParsed {
			return &PatternParsed{
				value: msg.Data["file"],
				msg:   msg,
			}
		},
	)
	*p = append(*p, *pattern)
	return p
}

// Reply type zero.PatternReplyMatched
func (p *Pattern) Reply(noCleanRedundantAt ...bool) *Pattern {
	noClean := false
	if len(noCleanRedundantAt) > 0 {
		noClean = noCleanRedundantAt[0]
	}
	pattern := NewPatternSegment(
		"reply", false, func(msg *message.Segment) *PatternParsed {
			return &PatternParsed{
				value: msg.Data["id"],
				msg:   msg,
			}
		}, !noClean,
	)
	*p = append(*p, *pattern)
	return p
}

// Any match any segment
func (p *Pattern) Any() *Pattern {
	pattern := NewPatternSegment(
		"any", false, func(msg *message.Segment) *PatternParsed {
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
			return &parsed
		},
	)
	*p = append(*p, *pattern)
	return p
}

func (s *PatternSegment) matchType(msg message.Segment) bool {
	return s.typ == msg.Type || s.typ == "any"
}
func mustMatchAllPatterns(pattern Pattern) bool {
	for _, p := range pattern {
		if p.optional {
			return false
		}
	}
	return true
}
func patternMatch(ctx *Ctx, pattern Pattern, msgs []message.Segment) bool {
	if mustMatchAllPatterns(pattern) && len(pattern) != len(msgs) {
		return false
	}
	patternState := make([]*PatternParsed, 0, 4)
	i := 0
	j := 0
	for i < len(pattern) {
		var parsed *PatternParsed
		if j < len(msgs) && pattern[i].matchType(msgs[j]) {
			parsed = pattern[i].parse(&msgs[j])
		} else {
			parsed = &PatternParsed{
				value: nil,
				msg:   nil,
			}
		}
		if j >= len(msgs) || !pattern[i].matchType(msgs[j]) || parsed.value == nil {
			if pattern[i].optional {
				patternState = append(patternState, parsed)
				i++
				continue
			}
			return false
		}
		patternState = append(patternState, parsed)
		i++
		j++
	}
	ctx.State[KeyPattern] = patternState
	return true
}
