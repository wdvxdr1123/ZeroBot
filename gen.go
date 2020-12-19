package zero

//go:generate go run github.com/a8m/syncmap -o "matcher_map_gen.go" -pkg zero -name matcherMap map[uint64]*Matcher

//go:generate go run github.com/a8m/syncmap -o "seq_map_gen.go" -pkg zero -name seqSyncMap map[uint64]chan<-APIResponse
