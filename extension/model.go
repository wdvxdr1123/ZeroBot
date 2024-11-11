package extension

import zero "github.com/wdvxdr1123/ZeroBot"

// PrefixModel is model of zero.PrefixRule
type PrefixModel struct {
	Prefix string `zero:"prefix"`
	Args   string `zero:"args"`
}

// SuffixModel is model of zero.SuffixRule
type SuffixModel struct {
	Suffix string `zero:"suffix"`
	Args   string `zero:"args"`
}

// CommandModel is model of zero.CommandRule
type CommandModel struct {
	Command string `zero:"command"`
	Args    string `zero:"args"`
}

// KeywordModel is model of zero.KeywordRule
type KeywordModel struct {
	Keyword string `zero:"keyword"`
}

// FullMatchModel is model of zero.FullMatchRule
type FullMatchModel struct {
	Matched string `zero:"matched"`
}

// RegexModel is model of zero.RegexRule
type RegexModel struct {
	Matched []string `zero:"regex_matched"`
}

// PatternModel is model of zero.PatternRule
type PatternModel struct {
	Matched []zero.PatternParsed `zero:"pattern_matched"`
}
