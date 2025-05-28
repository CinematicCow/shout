package scanner

type Scanner struct {
	Extensions   []string
	Directories  []string
	SkipPatterns []string
	outFile      string
}

func New(extensions, dirs, skip []string, outFile string) *Scanner {
	return &Scanner{
		Extensions:   extensions,
		Directories:  dirs,
		SkipPatterns: skip,
		outFile:      outFile,
	}
}
