// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// See plugin/doc.go for an explanation.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/kustomize/pkg/pgmconfig"
	"sigs.k8s.io/kustomize/pkg/plugins"
)

func main() {
	root := inputFileRoot()
	file, err := os.Open(root + ".go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	readToPackageMain(scanner, file.Name())

	w := NewWriter(root)
	defer w.close()

	// This particular phrasing is required.
	w.write(
		fmt.Sprintf(
			"// Code generated by pluginator on %s; DO NOT EDIT.",
			root))
	w.write("package builtin")

	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "//go:generate") {
			continue
		}
		if l == "var "+plugins.PluginSymbol+" plugin" {
			w.write("func New" + root + "Plugin() *" + root + "Plugin {")
			w.write("  return &" + root + "Plugin{}")
			w.write("}")
			continue
		}
		w.write(l)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func inputFileRoot() string {
	n := os.Getenv("GOFILE")
	if !strings.HasSuffix(n, ".go") {
		log.Fatalf("expecting .go suffix on %s", n)
	}
	return n[:len(n)-len(".go")]
}

func readToPackageMain(s *bufio.Scanner, f string) {
	gotMain := false
	for !gotMain && s.Scan() {
		gotMain = strings.HasPrefix(s.Text(), "package main")
	}
	if !gotMain {
		log.Fatalf("%s missing package main", f)
	}
}

type writer struct {
	root string
	f    *os.File
}

func NewWriter(r string) *writer {
	n := makeOutputFileName(r)
	f, err := os.Create(n)
	if err != nil {
		log.Fatalf("unable to create `%s`; %v", n, err)
	}
	return &writer{root: r, f: f}
}

func makeOutputFileName(root string) string {
	return filepath.Join(
		os.Getenv("GOPATH"),
		"src",
		pgmconfig.DomainName,
		pgmconfig.ProgramName,
		pgmconfig.PluginRoot,
		"builtin",
		root+".go")
}

func (w *writer) close() {
	fmt.Println("Generated " + w.root)
	w.f.Close()
}

func (w *writer) write(line string) {
	_, err := w.f.WriteString(w.filter(line) + "\n")
	if err != nil {
		log.Printf("Trouble writing: %s", line)
		log.Fatal(err)
	}
}

func (w *writer) filter(in string) string {
	if ok, newer := w.replace(in, "type plugin struct"); ok {
		return newer
	}
	if ok, newer := w.replace(in, "*plugin)"); ok {
		return newer
	}
	return in
}

// replace 'plugin' with 'FooPlugin' in context
// sensitive manner.
func (w *writer) replace(in, target string) (bool, string) {
	if !strings.Contains(in, target) {
		return false, ""
	}
	newer := strings.Replace(
		target, "plugin", w.root+"Plugin", 1)
	return true, strings.Replace(in, target, newer, 1)
}

