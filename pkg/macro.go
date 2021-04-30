package autostart

type Macro struct {
	Insert string //what the macro expands out to/replaces
	Deps   []string
	Prefix []string
	Pre    []string
	Post   []string
}

func injectMacros(cmdLine string, config Config) string {
	return cmdLine
}
