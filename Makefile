include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

SHELL := /bin/bash
export PATH := $(PWD)/bin:$(PATH)
PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test $(PKGS) generate install_deps

$(eval $(call golang-version-check,1.12))

test: generate $(PKGS)

$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)

generate:
	go generate

install_deps: golang-dep-vendor-deps
	$(call golang-dep-vendor)
	go build -o bin/mockgen    ./vendor/github.com/golang/mock/mockgen
