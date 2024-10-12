package zero

import (
	"github.com/wdvxdr1123/ZeroBot/message"
	"regexp"
	"strconv"
	"strings"
)

type Pattern []PatternSegment

func NewPattern() *Pattern {
	pattern := make(Pattern, 0)
	return &pattern
}

type PatternSegment struct {
	Type     string
	Optional bool
	Parse    func(msg *message.MessageSegment) *PatternParsed
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
	Msg   *message.MessageSegment
}

func (p PatternParsed) GetText() []string {
	return p.Value.([]string)
}
func (p PatternParsed) GetAt() string {
	return p.Value.(string)
}
func (p PatternParsed) GetImage() string {
	return p.Value.(string)
}
func (p PatternParsed) GetReply() string {
	return p.Value.(string)
}

// Text use regex to search a 'text' segment
func (p *Pattern) Text(regex string) *Pattern {
	re := regexp.MustCompile(regex)
	pattern := PatternSegment{
		Type: "text",
		Parse: func(msg *message.MessageSegment) *PatternParsed {
			s := msg.Data["text"]
			s = strings.Trim(s, " \n\r\t")
			matchString := re.MatchString(s)
			if matchString {
				return &PatternParsed{
					Valid: true,
					Value: re.FindStringSubmatch(s),
					Msg:   msg,
				}
			} else {
				return &PatternParsed{
					Valid: false,
					Value: nil,
					Msg:   nil,
				}
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
		Parse: func(msg *message.MessageSegment) *PatternParsed {
			qq, _ := strconv.ParseInt(msg.Data["qq"], 10, 64)
			if len(id) == 0 || len(id) == 1 && id[0] == msg.Data["qq"] {
				return &PatternParsed{
					Valid: true,
					Value: qq,
					Msg:   msg,
				}
			} else {
				return &PatternParsed{
					Valid: false,
					Value: nil,
					Msg:   nil,
				}
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
		Parse: func(msg *message.MessageSegment) *PatternParsed {
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
		Parse: func(msg *message.MessageSegment) *PatternParsed {
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
func containsOptional(pattern Pattern) bool {
	for _, p := range pattern {
		if p.Optional {
			return true
		}
	}
	return false
}
func patternMatch(ctx *Ctx, pattern Pattern, msgs []message.MessageSegment) bool {
	if !containsOptional(pattern) && len(pattern) != len(msgs) {
		return false
	}
	if _, ok := ctx.State[KEY_PATTERN]; !ok {
		ctx.State[KEY_PATTERN] = make([]*PatternParsed, 0, 1)
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
				ctx.State[KEY_PATTERN] = append(ctx.State[KEY_PATTERN].([]*PatternParsed), parsed)
				i++
				continue
			}
			return false
		}
		ctx.State[KEY_PATTERN] = append(ctx.State[KEY_PATTERN].([]*PatternParsed), parsed)
		i++
		j++
	}
	return true
}
