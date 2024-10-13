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
		for i := 1; i < len(ctx.Event.Message); i++ {
			if ctx.Event.Message[i-1].Type == "reply" && ctx.Event.Message[i].Type == "at" {
				// [reply][at]
				reply := ctx.GetMessage(ctx.Event.Message[i-1].Data["id"])
				if reply.MessageID.ID() == 0 || reply.Sender == nil || reply.Sender.ID == 0 {
					// failed to get history message
					msgs = append(msgs, ctx.Event.Message[i])
					continue
				}
				if strconv.FormatInt(reply.Sender.ID, 10) != ctx.Event.Message[i].Data["qq"] {
					// @ other user in reply
					msgs = append(msgs, ctx.Event.Message[i])
				}
			} else {
				msgs = append(msgs, ctx.Event.Message[i])
			}
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
	Type     string
	Optional bool
	Parse    func(msg *message.Segment) *PatternParsed
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
	Valid bool
	Value any
	Msg   *message.Segment
}

func (p PatternParsed) Text() []string {
	if !p.Valid {
		return nil
	}
	return p.Value.([]string)
}
func (p PatternParsed) At() string {
	if !p.Valid {
		return ""
	}
	return p.Value.(string)
}
func (p PatternParsed) Image() string {
	if !p.Valid {
		return ""
	}
	return p.Value.(string)
}
func (p PatternParsed) Reply() string {
	if !p.Valid {
		return ""
	}
	return p.Value.(string)
}

// Text use regex to search a 'text' segment
func (p *Pattern) Text(regex string) *Pattern {
	re := regexp.MustCompile(regex)
	pattern := PatternSegment{
		Type: "text",
		Parse: func(msg *message.Segment) *PatternParsed {
			s := msg.Data["text"]
			s = strings.Trim(s, " \n\r\t")
			matchString := re.MatchString(s)
			if matchString {
				return &PatternParsed{
					Valid: true,
					Value: re.FindStringSubmatch(s),
					Msg:   msg,
				}
			}

			return &PatternParsed{
				Valid: false,
				Value: nil,
				Msg:   nil,
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
		Type: "at",
		Parse: func(msg *message.Segment) *PatternParsed {
			if len(id) == 0 || len(id) == 1 && id[0] == msg.Data["qq"] {
				return &PatternParsed{
					Valid: true,
					Value: msg.Data["qq"],
					Msg:   msg,
				}
			}

			return &PatternParsed{
				Valid: false,
				Value: nil,
				Msg:   nil,
			}
		},
	}
	*p = append(*p, pattern)
	return p
}

// Image use regex to match an 'at' segment, if id is not empty, only match specific target
func (p *Pattern) Image() *Pattern {
	pattern := PatternSegment{
		Type: "image",
		Parse: func(msg *message.Segment) *PatternParsed {
			return &PatternParsed{
				Valid: true,
				Value: msg.Data["file"],
				Msg:   msg,
			}
		},
	}
	*p = append(*p, pattern)
	return p
}

// Reply type zero.PatternReplyMatched
func (p *Pattern) Reply() *Pattern {
	pattern := PatternSegment{
		Type: "reply",
		Parse: func(msg *message.Segment) *PatternParsed {
			return &PatternParsed{
				Valid: true,
				Value: msg.Data["id"],
				Msg:   msg,
			}
		},
	}
	*p = append(*p, pattern)
	return p
}
func mustMatchAllPatterns(pattern Pattern) bool {
	for _, p := range pattern {
		if !p.Optional {
			return false
		}
	}
	return true
}
func patternMatch(ctx *Ctx, pattern Pattern, msgs []message.Segment) bool {
	if mustMatchAllPatterns(pattern) && len(pattern) != len(msgs) {
		return false
	}
	if _, ok := ctx.State[KeyPattern]; !ok {
		ctx.State[KeyPattern] = make([]*PatternParsed, 0, 1)
	}
	i := 0
	j := 0
	for i < len(pattern) {
		var parsed *PatternParsed
		if j < len(msgs) && pattern[i].Type == (msgs[j].Type) {
			parsed = pattern[i].Parse(&msgs[j])
		} else {
			parsed = &PatternParsed{
				Valid: false,
				Value: nil,
				Msg:   nil,
			}
		}
		if j >= len(msgs) || pattern[i].Type != (msgs[j].Type) || !parsed.Valid {
			if pattern[i].Optional {
				ctx.State[KeyPattern] = append(ctx.State[KeyPattern].([]*PatternParsed), parsed)
				i++
				continue
			}
			return false
		}
		ctx.State[KeyPattern] = append(ctx.State[KeyPattern].([]*PatternParsed), parsed)
		i++
		j++
	}
	return true
}
