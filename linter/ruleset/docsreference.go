package ruleset

// DocsReference is a documentation reference link and/or source.
// Usually it's a web link.
type DocsReference string

const ToDoReference = DocsReference("TODO")

var docsReferenceMap = map[string]DocsReference{ // nolint:gochecknoglobals
	"CPY": DocsReference("https://docs.docker.com/engine/reference/builder/#copy"),
	"ENV": DocsReference("https://docs.docker.com/engine/reference/builder/#env"),
	"EXP": DocsReference("https://docs.docker.com/engine/reference/builder/#expose"),
	"FRM": DocsReference("https://docs.docker.com/engine/reference/builder/#from"),
	"RUN": DocsReference("https://docs.docker.com/engine/reference/builder/#run"),
	"STL": DocsReference("https://docs.docker.com/engine/reference/builder/#from"),
	"STS": DocsReference("https://docs.docker.com/engine/reference/builder/#from"),
	"USR": DocsReference("https://docs.docker.com/engine/reference/builder/#user"),
	"WKD": DocsReference("https://docs.docker.com/engine/reference/builder/#workdir"),
}
