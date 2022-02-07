package zero

// New 生成空组
func New() *MatcherGroup {
	return &MatcherGroup{}
}

// Engine is alias of MatcherGroup
type Engine = MatcherGroup

var defualtMatcherGroup = New()

// MatcherGroup is the pre_handler, post_handler manager
type MatcherGroup struct {
	handlers []Handler
	matchers []*Matcher
}

func (g *MatcherGroup) Use(middlewares ...Handler) {
	g.handlers = append(g.handlers, middlewares...)
}

func (g *MatcherGroup) Group(middlewares ...Handler) *MatcherGroup {
	group := &MatcherGroup{
		handlers: append(g.handlers, middlewares...),
	}
	return group
}

// Delete 移除该 MatcherGroup 注册的所有 Matchers
func (g *MatcherGroup) Delete() {
	for _, m := range g.matchers {
		m.Delete()
	}
}

func combine(h Handler, handlers []Handler) []Handler {
	if len(handlers) == 0 {
		return []Handler{h}
	}
	return append([]Handler{h}, handlers...)
}

// On 添加新的指定消息类型的匹配器(默认Engine)
func On(typ string, handlers ...Handler) *Matcher { return defualtMatcherGroup.On(typ, handlers...) }

// On 添加新的指定消息类型的匹配器
func (g *MatcherGroup) On(typ string, handlers ...Handler) *Matcher {
	matcher := &Matcher{
		Handlers:     combine(Type(typ), handlers),
		MatcherGroup: g,
	}
	g.matchers = append(g.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessage 消息触发器
func OnMessage(handlers ...Handler) *Matcher { return On("message", handlers...) }

// OnMessage 消息触发器
func (g *MatcherGroup) OnMessage(handlers ...Handler) *Matcher { return g.On("message", handlers...) }

// OnNotice 系统提示触发器
func OnNotice(handlers ...Handler) *Matcher { return On("notice", handlers...) }

// OnNotice 系统提示触发器
func (g *MatcherGroup) OnNotice(handlers ...Handler) *Matcher { return g.On("notice", handlers...) }

// OnRequest 请求消息触发器
func OnRequest(handlers ...Handler) *Matcher { return On("request", handlers...) }

// OnRequest 请求消息触发器
func (g *MatcherGroup) OnRequest(handlers ...Handler) *Matcher { return On("request", handlers...) }

// OnMetaEvent 元事件触发器
func OnMetaEvent(handlers ...Handler) *Matcher { return On("meta_event", handlers...) }

// OnMetaEvent 元事件触发器
func (g *MatcherGroup) OnMetaEvent(handlers ...Handler) *Matcher {
	return On("meta_event", handlers...)
}

// OnPrefix 前缀触发器
func OnPrefix(prefix string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnPrefix(prefix, handlers...)
}

// OnPrefix 前缀触发器
func (g *MatcherGroup) OnPrefix(prefix string, handlers ...Handler) *Matcher {
	return On("message", combine(PrefixRule(prefix), handlers)...)
}

// OnSuffix 后缀触发器
func OnSuffix(suffix string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnSuffix(suffix, handlers...)
}

// OnSuffix 后缀触发器
func (g *MatcherGroup) OnSuffix(suffix string, handlers ...Handler) *Matcher {
	return On("message", combine(SuffixRule(suffix), handlers)...)
}

// OnCommand 命令触发器
func OnCommand(commands string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnCommand(commands, handlers...)
}

// OnCommand 命令触发器
func (g *MatcherGroup) OnCommand(commands string, handlers ...Handler) *Matcher {
	return On("message", combine(CommandRule(commands), handlers)...)
}

// OnRegex 正则触发器
func OnRegex(regexPattern string, handlers ...Handler) *Matcher {
	return OnMessage(append([]Handler{RegexRule(regexPattern)}, handlers...)...)
}

// OnRegex 正则触发器
func (g *MatcherGroup) OnRegex(regexPattern string, handlers ...Handler) *Matcher {
	matcher := &Matcher{
		Handlers: append([]Handler{Type("message"), RegexRule(regexPattern)}, handlers...),
	}
	g.matchers = append(g.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnKeyword 关键词触发器
func OnKeyword(keyword string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnKeyword(keyword, handlers...)
}

// OnKeyword 关键词触发器
func (g *MatcherGroup) OnKeyword(keyword string, handlers ...Handler) *Matcher {
	return On("message", combine(KeywordRule(keyword), handlers)...)
}

// OnFullMatch 完全匹配触发器
func OnFullMatch(src string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnFullMatch(src, handlers...)
}

// OnFullMatch 完全匹配触发器
func (g *MatcherGroup) OnFullMatch(src string, handlers ...Handler) *Matcher {
	return On("message", combine(FullMatchRule(src), handlers)...)
}

// OnFullMatchGroup 完全匹配触发器组
func OnFullMatchGroup(src []string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnFullMatchGroup(src, handlers...)
}

// OnFullMatchGroup 完全匹配触发器组
func (g *MatcherGroup) OnFullMatchGroup(src []string, handlers ...Handler) *Matcher {
	return On("message", combine(FullMatchRule(src...), handlers)...)
}

// OnKeywordGroup 关键词触发器组
func OnKeywordGroup(keywords []string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnKeywordGroup(keywords, handlers...)
}

// OnKeywordGroup 关键词触发器组
func (g *MatcherGroup) OnKeywordGroup(keywords []string, handlers ...Handler) *Matcher {
	return On("message", combine(KeywordRule(keywords...), handlers)...)
}

// OnCommandGroup 命令触发器组
func OnCommandGroup(commands []string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnCommandGroup(commands, handlers...)
}

// OnCommandGroup 命令触发器组
func (g *MatcherGroup) OnCommandGroup(commands []string, handlers ...Handler) *Matcher {
	return g.On("message", append([]Handler{CommandRule(commands...)}, handlers...)...)
}

// OnPrefixGroup 前缀触发器组
func OnPrefixGroup(prefix []string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnPrefixGroup(prefix, handlers...)
}

// OnPrefixGroup 前缀触发器组
func (g *MatcherGroup) OnPrefixGroup(prefix []string, handlers ...Handler) *Matcher {
	return g.On("message", append([]Handler{PrefixRule(prefix...)}, handlers...)...)
}

// OnSuffixGroup 后缀触发器组
func OnSuffixGroup(suffix []string, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnSuffixGroup(suffix, handlers...)
}

// OnSuffixGroup 后缀触发器组
func (g *MatcherGroup) OnSuffixGroup(suffix []string, handlers ...Handler) *Matcher {
	return g.On("message", append([]Handler{SuffixRule(suffix...)}, handlers...)...)
}

// OnShell shell命令触发器
func OnShell(command string, model interface{}, handlers ...Handler) *Matcher {
	return defualtMatcherGroup.OnShell(command, model, handlers...)
}

// OnShell shell命令触发器
func (g *MatcherGroup) OnShell(command string, model interface{}, handlers ...Handler) *Matcher {
	return g.On("message", combine(ShellRule(command, model), handlers)...)
}
