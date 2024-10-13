package zero

import (
	"github.com/wdvxdr1123/ZeroBot/message"
	"regexp"
	"strconv"
	"strings"
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
		for i := 1; i < len(ctx.Event.Message); i++ {
			if ctx.Event.Message[i-1].Type == "reply" && ctx.Event.Message[i].Type == "at" {
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
	typ      string
	Optional bool
	parse    func(msg *message.Segment) *PatternParsed
}

// SetOptional set previous segment is optional, is v is empty, Optional will be true
// if Pattern is empty, panic
func (p *Pattern) SetOptional(v ...bool) *Pattern {
	if len(*p) == 0 {
		panic("pattern is empty")
	}
	if len(v) == 1 {
		(*p)[len(*p)-1].Optional = v[0]
	} else {
		(*p)[len(*p)-1].Optional = true
	}
	return p
}

// PatternParsed PatternRule parse result
type PatternParsed struct {
	value any
	msg   *message.Segment
}

func (p PatternParsed) Text() []string {
	if p.value == nil {
		return nil
	}
	return p.value.([]string)
}
func (p PatternParsed) At() string {
	if p.value == nil {
		return ""
	}
	return p.value.(string)
}
func (p PatternParsed) Image() string {
	if p.value == nil {
		return ""
	}
	return p.value.(string)
}
func (p PatternParsed) Reply() string {
	if p.value == nil {
		return ""
	}
	return p.value.(string)
}

// Text use regex to search a 'text' segment
func (p *Pattern) Text(regex string) *Pattern {
	re := regexp.MustCompile(regex)
	pattern := PatternSegment{
		typ: "text",
		parse: func(msg *message.Segment) *PatternParsed {
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
	}
	*p = append(*p, pattern)
	return p
}

// At use regex to match an 'at' segment, if id is not empty, only match specific target
func (p *Pattern) At(id ...string) *Pattern {
	if len(id) > 1 {
		panic("at pattern only support one id")
	}
	pattern := PatternSegment{
		typ: "at",
		parse: func(msg *message.Segment) *PatternParsed {
			if len(id) == 0 || len(id) == 1 && id[0] == msg.Data["qq"] {
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
	}
	*p = append(*p, pattern)
	return p
}

// Image use regex to match an 'at' segment, if id is not empty, only match specific target
func (p *Pattern) Image() *Pattern {
	pattern := PatternSegment{
		typ: "image",
		parse: func(msg *message.Segment) *PatternParsed {
			return &PatternParsed{
				value: msg.Data["file"],
				msg:   msg,
			}
		},
	}
	*p = append(*p, pattern)
	return p
}

// Reply type zero.PatternReplyMatched
func (p *Pattern) Reply() *Pattern {
	pattern := PatternSegment{
		typ: "reply",
		parse: func(msg *message.Segment) *PatternParsed {
			return &PatternParsed{
				value: msg.Data["id"],
				msg:   msg,
			}
		},
	}
	*p = append(*p, pattern)
	return p
}
func mustMatchAllPatterns(pattern Pattern) bool {
	for _, p := range pattern {
		if p.Optional {
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
		if j < len(msgs) && pattern[i].typ == (msgs[j].Type) {
			parsed = pattern[i].parse(&msgs[j])
		} else {
			parsed = &PatternParsed{
				value: nil,
				msg:   nil,
			}
		}
		if j >= len(msgs) || pattern[i].typ != (msgs[j].Type) || parsed.value == nil {
			if pattern[i].Optional {
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
