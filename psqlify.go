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
	bufInit := &bytes.Buffer{}
	bufFinal := &bytes.Buffer{}
	bufDataTable := &bytes.Buffer{}
	bufRelationTable := &bytes.Buffer{}

	for _, f := range targets {
		p.generateFiles(f, bufInit, bufFinal, bufDataTable, bufRelationTable)
	}

	return p.Artifacts()
}

// generateFiles generate psql files from buffer contents
func (p *PSQLModule) generateFiles(
	f pgs.File,
	bufInit *bytes.Buffer,
	bufFinal *bytes.Buffer,
	bufDataTable *bytes.Buffer,
	bufRelationTable *bytes.Buffer,
) {
	p.Push(f.Name().String())
	defer p.Pop()

	bufInit.Reset()
	bufFinal.Reset()
	bufDataTable.Reset()
	bufRelationTable.Reset()

	alter := false
	if ok, _ := p.Parameters().Bool("alter"); ok {
		alter = true
	}
	p.Debug("Param: alter=" + fmt.Sprintf("%v", alter))

	v := initPSQLVisitor(p, bufInit, bufFinal, bufDataTable, bufRelationTable, alter)
	p.CheckErr(pgs.Walk(v, f), "unable to generate psql")

	outInit := bufInit.String()
	outFinal := bufFinal.String()
	outTables := bufDataTable.String()
	outRelations := bufRelationTable.String()

	outName := f.InputPath().BaseName() + ".pb.psql"
	p.AddGeneratorFile(
		f.InputPath().SetBase("00_init_").String()+outName,
		outInit,
	)
	p.AddGeneratorFile(
		f.InputPath().SetBase("99_final_").String()+outName,
		outFinal,
	)
	p.AddGeneratorFile(
		f.InputPath().SetBase("10_tables_").String()+outName,
		outTables,
	)
	p.AddGeneratorFile(
		f.InputPath().SetBase("20_relations_").String()+outName,
		outRelations,
	)
}

// PSQLVisitor represent a visitor to walk the proto tree and analyse content
// (File, Messages, Fields, options)
type PSQLVisitor struct {
	pgs.Visitor
	pgs.DebuggerCommon
	initW          io.Writer
	finalW         io.Writer
	dataTableW     io.Writer
	relationTableW io.Writer
	definitions    []string
	alter          bool
}

func initPSQLVisitor(
	d pgs.DebuggerCommon,
	initW io.Writer,
	finalW io.Writer,
	dataTableW io.Writer,
	relationTableW io.Writer,
	alter bool,
) pgs.Visitor {
	v := PSQLVisitor{
		initW:          initW,
		finalW:         finalW,
		dataTableW:     dataTableW,
		relationTableW: relationTableW,
		DebuggerCommon: d,
		definitions:    []string{},
		alter:          alter,
	}
	v.Visitor = pgs.PassThroughVisitor(&v)
	return &v
}

func writeComment(w io.Writer, str string) {
	fmt.Fprintf(w, "-- %s\n", str)
}

func write(w io.Writer, str string) {
	fmt.Fprintf(w, "%s", str)
}

func writeln(w io.Writer, str string) {
	fmt.Fprintf(w, "%s\n", str)
}

// createTable start a statement to create a table
func (v *PSQLVisitor) createTable(w io.Writer, str string) {
	write(w, "CREATE TABLE IF NOT EXISTS "+str+" (")

	var strEnd string
	if v.alter {
		strEnd = ");"
	} else {
		strEnd = ""
	}

	writeln(w, strEnd)
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
func (v *PSQLVisitor) writeDefinition(w io.Writer) {
	v.Debug("writeDefinition: defs = %v+", v.definitions)
	if v.alter {
		fmt.Fprintf(w, "%s\n", strings.Join(v.definitions, "\n"))
	} else {
		fmt.Fprintf(w, "\t%s\n", strings.Join(v.definitions, ",\n\t"))
	}
	v.definitions = []string{}
}

// VisitFile prepare a .psql from a proto one
// For each messages, call VisitMessage
func (v *PSQLVisitor) VisitFile(f pgs.File) (pgs.Visitor, error) {
	v.Debug("pssql: Processing file " + f.Name().String())
	writeComment(v.initW, "File: "+f.Name().String())
	write(v.initW, "")

	initializations := make([]string, 0)
	finalizations := make([]string, 0)

	if ok, err := f.Extension(psql.E_Initialization, &initializations); ok && err == nil {
		for _, init := range initializations {
			writeln(v.initW, init)
		}
	}

	for _, message := range f.AllMessages() {
		pgs.Walk(v, message)
	}

	if ok, err := f.Extension(psql.E_Finalization, &finalizations); ok && err == nil {
		for _, final := range finalizations {
			writeln(v.finalW, final)
		}
	}

	return nil, nil
}

// VisitMessage extract psql related options of a messages and generate associated statements
// For each fields, call VisitField
func (v *PSQLVisitor) VisitMessage(m pgs.Message) (pgs.Visitor, error) {
	var disabled bool
	var tableType psql.TableType
	var w *io.Writer

	if ok, err := m.Extension(psql.E_Disabled, &disabled); ok && err == nil && disabled {
		v.Logf("Generation disabled for message " + m.Name().String())
		return nil, nil
	}

	ok, err := m.Extension(psql.E_TableType, &tableType)
	if err != nil {
		v.Logf(err.Error())
		return nil, nil
	}

	if !ok {
		v.Logf("Unable to find an extension tableType equal to DATA or RELATION. Skipping message: " + m.Name().String())
		return nil, nil
	}

	switch tableType {
	case psql.TableType_DATA:
		v.Debug("Table Type: DATA")
		w = &v.dataTableW
	case psql.TableType_RELATION:
		v.Debug("Table Type: RELATION")
		w = &v.relationTableW
	default:
		w = &v.dataTableW
	}

	writeComment(*w, "File: "+m.File().Name().String())
	v.createTable(*w, m.Name().String())

	prefixes := make([]string, 0)
	suffixes := make([]string, 0)
	constraints := make([]string, 0)

	if ok, err := m.Extension(psql.E_Prefix, &prefixes); ok && err == nil {
		for _, prefix := range prefixes {
			v.appendDefinition(prefix)
		}
	}

	for _, field := range m.Fields() {
		pgs.Walk(v, field)
	}

	if ok, err := m.Extension(psql.E_Constraint, &constraints); ok && err == nil {
		for _, constraint := range constraints {
			v.appendConstraint(m.Name().String(), constraint)
		}
	}

	if ok, err := m.Extension(psql.E_Suffix, &suffixes); ok && err == nil {
		for _, suffix := range suffixes {
			v.appendDefinition(suffix)
		}
	}
	v.writeDefinition(*w)

	// close the CREATE TABLE statement properly if we are not in "alter" mode
	if !v.alter {
		writeln(*w, ");")
	}

	return nil, nil
}

// VisitField extract psql related options of a field and generate associated statements
func (v *PSQLVisitor) VisitField(f pgs.Field) (pgs.Visitor, error) {
	var column string

	ok, err := f.Extension(psql.E_Column, &column)
	if ok && err == nil {
		v.appendColumn(f.Message().Name().String(), f.Name().String()+" "+column)

	}
	return nil, nil
}
