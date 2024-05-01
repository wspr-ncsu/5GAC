package instrumentation

import (
	"github.com/wspr-ncsu/5GAC/gocode/Shared/accesspattern"
)

type Parser interface {
	Parse(file string)
}

type Inspecter interface {
	Inspect()
}

type GoFile accesspattern.GoFile
