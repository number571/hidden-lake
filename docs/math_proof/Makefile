ENV=$(HOME)/.elan/env
.PHONY: default install build
default: install build
install:
	source $(ENV) && lake env lean MathProof.lean
	source $(ENV) && lake exe cache get
build:
	source $(ENV) && lake build
