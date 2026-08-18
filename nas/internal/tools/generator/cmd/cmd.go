package main

import "github.com/acore2026/free6gc-lib/nas/internal/tools/generator"

func main() {
	generator.ParseSpecs()

	generator.GenerateNasMessage()

	generator.GenerateNasEncDec()

	generator.GenerateTestLarge()
}
