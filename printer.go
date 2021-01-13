package main

import (
	"fmt"
	"io"
	"log"

	"bytes"

	psql "github.com/intrinsec/protoc-gen-psql/psql"
	pgs "github.com/lyft/protoc-gen-star"
)

type PrinterModule struct {
	*pgs.ModuleBase
}

func ASTPrinter() *PrinterModule { return &PrinterModule{ModuleBase: &pgs.ModuleBase{}} }

func (p *PrinterModule) Name() string { return "psql" }

func (p *PrinterModule) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	buf := &bytes.Buffer{}

	for _, f := range targets {
		p.printFile(f, buf)
	}

	return p.Artifacts()
}

func (p *PrinterModule) printFile(f pgs.File, buf *bytes.Buffer) {
	p.Push(f.Name().String())
	defer p.Pop()

	buf.Reset()
	v := initPrintVisitor(buf)
	p.CheckErr(pgs.Walk(v, f), "unable to generate psql")

	out := buf.String()

	p.AddGeneratorFile(
		f.InputPath().SetExt(".psql").String(),
		out,
	)
}

type PrinterVisitor struct {
	pgs.Visitor
	w io.Writer
}

func initPrintVisitor(w io.Writer) pgs.Visitor {
	v := PrinterVisitor{
		w: w,
	}
	v.Visitor = pgs.PassThroughVisitor(&v)
	return v
}

func (v PrinterVisitor) writeComment(str string) {
	fmt.Fprintf(v.w, "-- %s\n", str)
}

func (v PrinterVisitor) write(str string) {
	fmt.Fprintf(v.w, "%s\n", str)
}

func (v PrinterVisitor) writeIndented(str string) {
	fmt.Fprintf(v.w, "\t%s,\n", str)
}

func (v PrinterVisitor) VisitFile(f pgs.File) (pgs.Visitor, error) {
	log.Println("pssql: Processing file " + f.Name().String())
	v.writeComment("File: " + f.Name().String())
	v.write("")

	initializations := make([]string, 0)
	finalizations := make([]string, 0)

	if ok, err := f.Extension(psql.E_Initializations, &initializations); ok && err == nil {
		for _, init := range initializations {
			v.write(init)
		}
	}

	v.write("")

	for _, message := range f.AllMessages() {
		pgs.Walk(v, message)
	}

	v.write("")

	if ok, err := f.Extension(psql.E_Finalizations, &finalizations); ok && err == nil {
		for _, finit := range finalizations {
			v.write(finit)
		}
	}

	return nil, nil
}

func (v PrinterVisitor) VisitMessage(m pgs.Message) (pgs.Visitor, error) {

	var disabled bool

	if ok, err := m.Extension(psql.E_Disabled, &disabled); ok && err == nil && disabled {
		log.Println("pssql: Generation disabled for message " + m.Name().String())
		return nil, nil
	}

	prefixes := make([]string, 0)
	suffixes := make([]string, 0)

	v.write("CREATE TABLE IF NOT EXISTS " + m.Name().String() + " (")

	if ok, err := m.Extension(psql.E_Prefixes, &prefixes); ok && err == nil {
		for _, prefix := range prefixes {
			v.writeIndented(prefix)
		}
	}

	for _, field := range m.Fields() {
		pgs.Walk(v, field)
	}

	if ok, err := m.Extension(psql.E_Suffixes, &suffixes); ok && err == nil {
		for _, suffix := range suffixes {
			v.writeIndented(suffix)
		}
	}

	v.write(");")

	return nil, nil
}

func (v PrinterVisitor) VisitField(f pgs.Field) (pgs.Visitor, error) {

	var column string

	ok, err := f.Extension(psql.E_Column, &column)
	if ok && err == nil {
		v.writeIndented(f.Name().String() + " " + column)
	} else {
		v.writeComment("Ignored Field: " + f.Name().String())
	}
	return nil, nil
}
