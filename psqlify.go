package main

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"bytes"

	psql "github.com/intrinsec/protoc-gen-psql/psql"
	pgs "github.com/lyft/protoc-gen-star"
)

var (
	//go:embed _templates/auto_fill_on_update.tpl.psql
	templateAutoFillOnUpdate string

	//go:embed _templates/cascade_update.tpl.psql
	templateCascadeUpdate string

	//go:embed _templates/create_table.tpl.psql
	templateCreateTable string

	//go:embed _templates/alter_table.tpl.psql
	templateAlterTable string

	//go:embed _templates/generic.tpl.psql
	templateGeneric string
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
	columns        []string
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
		columns:        []string{},
		alter:          alter,
	}
	v.Visitor = pgs.PassThroughVisitor(&v)
	return &v
}

// writeAutoFillUpdate write auto fill function and trigger to final psql file
func (v *PSQLVisitor) writeAutoFillUpdate(t string, field string, value string) {
	data := struct {
		FunctionName string
		TriggerName  string
		Table        string
		Field        string
		Value        string
	}{
		FunctionName: fmt.Sprintf("fn_auto_fill_%s_on_%s_update", field, strings.ToLower(t)),
		TriggerName:  fmt.Sprintf("tg_auto_fill_%s_on_%s_update", field, strings.ToLower(t)),
		Table:        t,
		Field:        field,
		Value:        value,
	}
	generateFromTemplate(templateAutoFillOnUpdate, data.TriggerName, data, v.finalW)
}

// addAutoFillUpdate write auto fill function and trigger to final psql file
func (v *PSQLVisitor) writeCascadeUpdate(t string, field string, cascade_update *psql.CascadeUpdate) {
	for _, updates := range cascade_update.GetUpdates() {
		data := struct {
			FunctionName  string
			TriggerName   string
			Table         string
			Field         string
			TableToUpdate string
			Key           string
			FieldToUpdate string
			Value         string
		}{
			FunctionName:  fmt.Sprintf("fn_%s_cascade_update_%s_on_%s", strings.ToLower(t), updates.Field, strings.ToLower(cascade_update.Table)),
			TriggerName:   fmt.Sprintf("tg_%s_cascade_update_%s_on_%s", strings.ToLower(t), updates.Field, strings.ToLower(cascade_update.Table)),
			Table:         t,
			Field:         field,
			TableToUpdate: cascade_update.Table,
			Key:           cascade_update.Key,
			FieldToUpdate: updates.Field,
			Value:         updates.Value,
		}
		generateFromTemplate(templateCascadeUpdate, data.TriggerName, data, v.finalW)
	}
}

func generateFromTemplate(templateText string, templateName string, data interface{}, writer io.Writer) {
	tmpl, err := template.New(templateName).Funcs(template.FuncMap{"StringsJoin": strings.Join}).Parse(templateText)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(writer, data)
	if err != nil {
		panic(err)
	}
}

func AppendSlices(slices [][]string) []string {
	var tmp []string
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}

// VisitFile prepare a .psql from a proto one
// For each messages, call VisitMessage
func (v *PSQLVisitor) VisitFile(f pgs.File) (pgs.Visitor, error) {
	v.Debugf("pssql: Processing file %s", f.Name().String())

	initializations := make([]string, 0)
	finalizations := make([]string, 0)

	if ok, err := f.Extension(psql.E_Initialization, &initializations); ok && err != nil {
		v.Logf("Error can't retrieve initialization extensions for file %s with error: %s", f.Name().String(), err)
	}

	if ok, err := f.Extension(psql.E_Finalization, &finalizations); ok && err != nil {
		v.Logf("Error can't retrieve finalization extensions for file %s with error: %s", f.Name().String(), err)
	}

	data := struct {
		FileName string
		Queries  []string
	}{
		FileName: f.Name().String(),
		Queries:  initializations,
	}
	generateFromTemplate(templateGeneric, "initialization "+f.Name().String(), data, v.initW)

	data.Queries = finalizations
	generateFromTemplate(templateGeneric, "finalization"+f.Name().String(), data, v.finalW)

	return initPSQLVisitor(v, v.initW, v.finalW, v.dataTableW, v.relationTableW, v.alter), nil
}

// VisitMessage extract psql related options of a messages and generate associated statements
// For each fields, call VisitField
func (v *PSQLVisitor) VisitMessage(m pgs.Message) (pgs.Visitor, error) {
	var disabled bool
	var tableType psql.TableType
	var w *io.Writer

	if ok, err := m.Extension(psql.E_Disabled, &disabled); ok && err == nil && disabled {
		v.Logf("Generation disabled for message %s", m.Name().String())
		return nil, nil
	}

	ok, err := m.Extension(psql.E_TableType, &tableType)
	if err != nil {
		v.Logf(err.Error())
		return nil, nil
	}

	if !ok {
		v.Logf("Unable to find an extension tableType equal to DATA or RELATION. Skipping message: %s", m.Name().String())
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

	prefixes := make([]string, 0)
	suffixes := make([]string, 0)
	constraints := make([]string, 0)

	if ok, err := m.Extension(psql.E_Prefix, &prefixes); ok && err != nil {
		v.Logf("Error can't retrieve prefix extensions for message %s with error: %s", m.Name().String(), err)
	}

	if ok, err := m.Extension(psql.E_Constraint, &constraints); ok && err != nil {
		v.Logf("Error can't retrieve constraint extensions for message %s with error: %s", m.Name().String(), err)
	}

	if ok, err := m.Extension(psql.E_Suffix, &suffixes); ok && err != nil {
		v.Logf("Error can't retrieve suffix extensions for message %s with error: %s", m.Name().String(), err)
	}

	for _, field := range m.Fields() {
		pgs.Walk(v, field)
	}

	var templateText string
	if v.alter {
		templateText = templateAlterTable
	} else {
		templateText = templateCreateTable
	}

	definitions := AppendSlices([][]string{prefixes, v.columns, constraints, suffixes})

	data := struct {
		FileName    string
		Name        string
		Definitions []string
		Prefixes    []string
		Columns     []string
		Constraints []string
		Suffixes    []string
	}{
		FileName:    m.File().Name().String(),
		Name:        m.Name().String(),
		Definitions: definitions,
		Prefixes:    prefixes,
		Columns:     v.columns,
		Constraints: constraints,
		Suffixes:    suffixes,
	}
	generateFromTemplate(templateText, "dataTable "+m.Name().String(), data, *w)
	v.columns = []string{}

	return nil, nil
}

// VisitField extract psql related options of a field and generate associated statements
func (v *PSQLVisitor) VisitField(f pgs.Field) (pgs.Visitor, error) {
	var column string
	var auto_fill_value string
	var cascade_update psql.CascadeUpdate

	ok, err := f.Extension(psql.E_Column, &column)
	if ok && err == nil {
		v.columns = append(v.columns, f.Name().String()+" "+column)
	}

	ok, err = f.Extension(psql.E_AutoFillOnUpdate, &auto_fill_value)
	if ok && err == nil {
		v.writeAutoFillUpdate(f.Message().Name().String(), f.Name().String(), auto_fill_value)
	}

	ok, err = f.Extension(psql.E_CascadeUpdate, &cascade_update)
	if ok && err == nil {
		v.writeCascadeUpdate(f.Message().Name().String(), f.Name().String(), &cascade_update)
	}

	return nil, nil
}
