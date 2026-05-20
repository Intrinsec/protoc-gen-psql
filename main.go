package main

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

func main() {
	pgs.
		Init(pgs.DebugEnv("DEBUG_PGV")).
		RegisterModule(PSQLify()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
