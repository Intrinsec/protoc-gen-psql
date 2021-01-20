package main

import (
	"fmt"
	"io"
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

	alter := false
	if ok, _ := p.Parameters().Bool("alter"); ok {
		alter = true
	}
	p.Debug("Param: alter=" + fmt.Sprintf("%v", alter))
	v := initPSQLVisitor(p, buf, alter)
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
	pgs.DebuggerCommon
	w           io.Writer
	definitions []string
	alter       bool
}

func initPSQLVisitor(d pgs.DebuggerCommon, w io.Writer, alter bool) pgs.Visitor {
	v := PSQLVisitor{
		w:              w,
		DebuggerCommon: d,
		definitions:    []string{},
		alter:          alter,
	}
	v.Visitor = pgs.PassThroughVisitor(&v)
	return &v
}

func (v *PSQLVisitor) writeComment(str string) {
	fmt.Fprintf(v.w, "-- %s\n", str)
}

func (v *PSQLVisitor) write(str string) {
	fmt.Fprintf(v.w, "%s", str)
}

func (v *PSQLVisitor) writeln(str string) {
	fmt.Fprintf(v.w, "%s\n", str)
}

// createTable start a statement to create a table
func (v *PSQLVisitor) createTable(str string) {
	v.write("CREATE TABLE IF NOT EXISTS " + str + " (")

	var strEnd string
	if v.alter {
		strEnd = ");"
	} else {
		strEnd = ""
	}

	v.writeln(strEnd)
}

// appendDefinition append a full independent SQL statement without any formatting
func (v *PSQLVisitor) appendDefinition(str string) {
	v.definitions = append(v.definitions, str)
	v.Debug("appendDefinition: defs = %v+", v.definitions)
}

// appendColumn append a column col to a table t
func (v *PSQLVisitor) appendColumn(t string, col string) {
	var def string
	if v.alter {
		def = "ALTER TABLE " + t + " ADD COLUMN IF NOT EXISTS " + col + ";\n"
	} else {
		def = col
	}

	v.definitions = append(v.definitions, def)
}

// appendConstraint append a constraint str to a table t
func (v *PSQLVisitor) appendConstraint(t string, str string) {
	var cs string
	if v.alter {
		cs = "DO $$\n" +
			"BEGIN\n" +
			"ALTER TABLE " + t + " ADD " + str + ";\n" +
			"EXCEPTION WHEN duplicate_object THEN RAISE NOTICE '%, skipping', SQLERRM USING ERRCODE = SQLSTATE;\n" +
			"END\n" +
			"$$;\n"
	} else {
		cs = str
	}
	v.definitions = append(v.definitions, cs)
}

// writeDefinition write definitions to a file then clear definitions
func (v *PSQLVisitor) writeDefinition() {
	v.Debug("writeDefinition: defs = %v+", v.definitions)
	if v.alter {
		fmt.Fprintf(v.w, "%s\n", strings.Join(v.definitions, "\n"))
	} else {
		fmt.Fprintf(v.w, "\t%s\n", strings.Join(v.definitions, ",\n\t"))
	}
	v.definitions = []string{}
}

// VisitFile prepare a .psql from a proto one
// For each messages, call VisitMessage
func (v *PSQLVisitor) VisitFile(f pgs.File) (pgs.Visitor, error) {
	v.Debug("pssql: Processing file " + f.Name().String())
	v.writeComment("File: " + f.Name().String())
	v.write("")

	initializations := make([]string, 0)
	finalizations := make([]string, 0)

	if ok, err := f.Extension(psql.E_Initializations, &initializations); ok && err == nil {
		for _, init := range initializations {
			v.writeln(init)
		}
	}

	v.writeln("")

	for _, message := range f.AllMessages() {
		pgs.Walk(v, message)
	}

	v.writeln("")

	if ok, err := f.Extension(psql.E_Finalizations, &finalizations); ok && err == nil {
		for _, finit := range finalizations {
			v.writeln(finit)
		}
	}

	return nil, nil
}

// VisitMessage extract psql related options of a messages and generate associated statements
// For each fields, call VisitField
func (v *PSQLVisitor) VisitMessage(m pgs.Message) (pgs.Visitor, error) {
	var disabled bool

	if ok, err := m.Extension(psql.E_Disabled, &disabled); ok && err == nil && disabled {
		v.Logf("pssql: Generation disabled for message " + m.Name().String())
		return nil, nil
	}

	prefixes := make([]string, 0)
	suffixes := make([]string, 0)
	constraints := make([]string, 0)

	v.createTable(m.Name().String())

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

	// close the CREATE TABLE statement properly if we are not in "alter" mode
	if !v.alter {
		v.writeln(");")
	}

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
