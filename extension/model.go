package extension

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
