package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"bytes"

	psql "github.com/intrinsec/protoc-gen-psql/psql"
	pgs "github.com/lyft/protoc-gen-star"
)

// PSQLModule implement a custom protoc-gen-star module
type PSQLModule struct {
	*pgs.ModuleBase
}

// PSQLify returns and initialized PSQLify module
func PSQLify() *PSQLModule {
	return &PSQLModule{ModuleBase: &pgs.ModuleBase{}}
}

// Name define the name of the protoc module
func (p *PSQLModule) Name() string { return "psql" }

// Execute generates PSQL files from received proto ones
func (p *PSQLModule) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	buf := &bytes.Buffer{}

	for _, f := range targets {
		p.generateFile(f, buf)
	}

	return p.Artifacts()
}

// generateFile generate a psql file from a buffer content
func (p *PSQLModule) generateFile(f pgs.File, buf *bytes.Buffer) {
	p.Push(f.Name().String())
	defer p.Pop()

	buf.Reset()
	v := initPSQLVisitor(buf)
	p.CheckErr(pgs.Walk(v, f), "unable to generate psql")

	out := buf.String()

	p.AddGeneratorFile(
		f.InputPath().SetExt(".psql").String(),
		out,
	)
}

// PSQLVisitor represent a visitor to walk the proto tree and analyse content
// (File, Messages, Fields, options)
type PSQLVisitor struct {
	pgs.Visitor
	w           io.Writer
	definitions []string
}

func initPSQLVisitor(w io.Writer) pgs.Visitor {
	v := PSQLVisitor{
		w:           w,
		definitions: []string{},
	}
	v.Visitor = pgs.PassThroughVisitor(&v)
	return &v
}

func (v *PSQLVisitor) writeComment(str string) {
	fmt.Fprintf(v.w, "-- %s\n", str)
}

func (v *PSQLVisitor) write(str string) {
	fmt.Fprintf(v.w, "%s\n", str)
}

// appendDefinition append a full independent SQL statement without any formatting
func (v *PSQLVisitor) appendDefinition(str string) {
	v.definitions = append(v.definitions, str)
	log.Printf("appendDefinition: defs = %v+", v.definitions)
}

// appendColumn append a column col to a table t
func (v *PSQLVisitor) appendColumn(t string, col string) {
	def := "ALTER TABLE " + t + " ADD COLUMN IF NOT EXISTS " + col + ";\n"

	v.definitions = append(v.definitions, def)
}

// appendConstraint append a constraint str to a table t
func (v *PSQLVisitor) appendConstraint(t string, str string) {
	cs := "DO $$\n" +
		"BEGIN\n" +
		"ALTER TABLE " + t + " ADD " + str + ";\n" +
		"EXCEPTION WHEN duplicate_object THEN RAISE NOTICE '%, skipping', SQLERRM USING ERRCODE = SQLSTATE;\n" +
		"END\n" +
		"$$;\n"
	v.definitions = append(v.definitions, cs)
}

// writeDefinition write definitions to a file then clear definitions
func (v *PSQLVisitor) writeDefinition() {
	log.Printf("writeDefinition: defs = %v+", v.definitions)
	fmt.Fprintf(v.w, "%s\n", strings.Join(v.definitions, "\n"))
	v.definitions = []string{}
}

// VisitFile prepare a .psql from a proto one
// For each messages, call VisitMessage
func (v *PSQLVisitor) VisitFile(f pgs.File) (pgs.Visitor, error) {
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

// VisitMessage extract psql related options of a messages and generate associated statements
// For each fields, call VisitField
func (v *PSQLVisitor) VisitMessage(m pgs.Message) (pgs.Visitor, error) {

	var disabled bool

	if ok, err := m.Extension(psql.E_Disabled, &disabled); ok && err == nil && disabled {
		log.Println("pssql: Generation disabled for message " + m.Name().String())
		return nil, nil
	}

	prefixes := make([]string, 0)
	suffixes := make([]string, 0)
	constraints := make([]string, 0)

	v.write("CREATE TABLE IF NOT EXISTS " + m.Name().String() + " ();\n")

	if ok, err := m.Extension(psql.E_Prefixes, &prefixes); ok && err == nil {
		for _, prefix := range prefixes {
			v.appendDefinition(prefix)
		}
	}

	for _, field := range m.Fields() {
		pgs.Walk(v, field)
	}

	if ok, err := m.Extension(psql.E_Constraints, &constraints); ok && err == nil {
		for _, constraint := range constraints {
			v.appendConstraint(m.Name().String(), constraint)
		}
	}

	if ok, err := m.Extension(psql.E_Suffixes, &suffixes); ok && err == nil {
		for _, suffix := range suffixes {
			v.appendDefinition(suffix)
		}
	}
	v.writeDefinition()

	return nil, nil
}

// VisitField extract psql related options of a field and generate associated statements
func (v *PSQLVisitor) VisitField(f pgs.Field) (pgs.Visitor, error) {

	var column string

	ok, err := f.Extension(psql.E_Column, &column)
	if ok && err == nil {
		v.appendColumn(f.Message().Name().String(), f.Name().String()+" "+column)

	} else {
		v.writeComment("Ignored Field: " + f.Name().String())
	}
	return nil, nil
}
