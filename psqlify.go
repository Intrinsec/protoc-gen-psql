package main

import (
	_ "embed"
	"fmt"
	"hash/adler32"
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

	//go:embed _templates/relay_cascade_update.tpl.psql
	templateRelayCascadeUpdate string

	//go:embed _templates/cascade_update_on_related_table.tpl.psql
	templateCascadeUpdateOnRelatedTable string

	//go:embed _templates/create_table.tpl.psql
	templateCreateTable string

	//go:embed _templates/alter_table.tpl.psql
	templateAlterTable string
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

	filePath := f.InputPath().String()
	outName := f.InputPath().BaseName() + ".pb.psql"

	if outInit, count := getStringBufferWithHeader(bufInit, filePath); count != 0 {
		p.AddGeneratorFile(
			f.InputPath().SetBase("00_init_").String()+outName,
			outInit,
		)
	}
	if outFinal, count := getStringBufferWithHeader(bufFinal, filePath); count != 0 {
		p.AddGeneratorFile(
			f.InputPath().SetBase("99_final_").String()+outName,
			outFinal,
		)
	}
	if outTables, count := getStringBufferWithHeader(bufDataTable, filePath); count != 0 {
		p.AddGeneratorFile(
			f.InputPath().SetBase("10_tables_").String()+outName,
			outTables,
		)
	}
	if outRelations, count := getStringBufferWithHeader(bufRelationTable, filePath); count != 0 {
		p.AddGeneratorFile(
			f.InputPath().SetBase("20_relations_").String()+outName,
			outRelations,
		)
	}
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
func (v *PSQLVisitor) writeAutoFillUpdate(table string, field string, value string) {
	fnName, tgName, _, createFnName := generateCascadeIdentifierNames("auto_fill", table, field)

	data := struct {
		FunctionName       string
		TriggerName        string
		CreateFunctionName string
		Table              string
		Field              string
		Value              string
	}{
		FunctionName:       fnName,
		TriggerName:        tgName,
		CreateFunctionName: createFnName,
		Table:              table,
		Field:              field,
		Value:              value,
	}
	generateFromTemplate(templateAutoFillOnUpdate, data.TriggerName, data, v.finalW)
}

func (v *PSQLVisitor) writeRelayCascadeUpdate(relationTable string, relayCascadeUpdates []*psql.RelayCascadeUpdate) {

	for _, relayCascadeUpdate := range relayCascadeUpdates {
		for _, destination := range relayCascadeUpdate.Destinations {
			fnName, tgName, tgDelName, createFnName := generateCascadeIdentifierNames("relay_cascade", relationTable, relayCascadeUpdate.SourceForeignKey, destination.ForeignKey)

			data := struct {
				FunctionName          string
				TriggerName           string
				TriggerDelName        string
				CreateFunctionName    string
				RelationTable         string
				FieldToUpdate         string
				SourceForeignKey      string
				DestinationForeignKey string
				Value                 string
			}{
				FunctionName:          fnName,
				TriggerName:           tgName,
				TriggerDelName:        tgDelName,
				CreateFunctionName:    createFnName,
				RelationTable:         relationTable,
				FieldToUpdate:         destination.Field,
				SourceForeignKey:      relayCascadeUpdate.SourceForeignKey,
				DestinationForeignKey: destination.ForeignKey,
				Value:                 destination.Value,
			}
			generateFromTemplate(templateRelayCascadeUpdate, data.TriggerName, data, v.finalW)
		}
	}
}

func (v *PSQLVisitor) writeCascadeUpdateOnRelatedTable(relationTable string, foreignKey string, cascadeUpdateOnRelatedTables []*psql.CascadeUpdateOnRelatedTable) {
	fnName, tgName, tgDelName, createFnName := generateCascadeIdentifierNames("cascade_related", relationTable, foreignKey)

	data := struct {
		FunctionName       string
		TriggerName        string
		TriggerDelName     string
		CreateFunctionName string
		RelationTable      string
		ForeignKey         string
		Updates            []*psql.CascadeUpdateOnRelatedTable
	}{
		FunctionName:       fnName,
		TriggerName:        tgName,
		TriggerDelName:     tgDelName,
		CreateFunctionName: createFnName,
		RelationTable:      relationTable,
		ForeignKey:         foreignKey,
		Updates:            cascadeUpdateOnRelatedTables,
	}
	generateFromTemplate(templateCascadeUpdateOnRelatedTable, data.TriggerName, data, v.finalW)
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

func appendSlices(slices ...[]string) []string {
	var tmp []string
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}

func generateCascadeIdentifierNames(name string, parameters ...string) (fnName string, tgName string, tgDelName string, createFnName string) {

	identifierNames := map[string]string{
		"fnName":       "fn_",
		"tgName":       "tg_",
		"tgDelName":    "tg_del_",
		"fnCreateName": "fn_create_",
	}
	maxPrefixLen := 0

	for _, v := range identifierNames {
		if len(v) > maxPrefixLen {
			maxPrefixLen = len(v)
		}
	}

	baseIdentifier, err := generateIdentifierName(name, maxPrefixLen, parameters...)
	if err != nil {
		panic(err)
	}

	fnName = identifierNames["fnName"] + baseIdentifier
	tgName = identifierNames["tgName"] + baseIdentifier
	tgDelName = identifierNames["tgDelName"] + baseIdentifier
	createFnName = identifierNames["fnCreateName"] + baseIdentifier

	return
}

// generate a unique and valid postgresql identifier name which can be used as function or trigger names
func generateIdentifierName(name string, prefixSize int, parameters ...string) (string, error) {
	const MAX_PSQL_IDENTIFIER_SIZE = 63
	const CHECKSUM_SIZE = 8

	identifier := name

	if prefixSize > 10 {
		return "", fmt.Errorf("given prefix size %d is longer than 10 characters", prefixSize)
	}

	totalParametersSize := (MAX_PSQL_IDENTIFIER_SIZE -
		len(name) -
		prefixSize -
		(len(parameters) + 1) - // number of parameter and checksum '_' separators
		CHECKSUM_SIZE)

	// iterate over parameters instead of parameterSizeMap to keep a consistant order (hashmap is unordered).
	parameterSizeMap := dispatchMaxSizeByParameter(totalParametersSize, parameters...)
	for _, parameter := range parameters {
		size := parameterSizeMap[parameter]
		identifier += fmt.Sprintf("_%s", parameter[:size])
	}

	// compute the checksum over all non-truncated parameters to avoid collision
	checksumData := []byte(strings.Join(parameters, "-"))
	identifier += fmt.Sprintf("_%x", adler32.Checksum(checksumData))

	identifier = strings.ToLower(identifier)

	if prefixSize+len(identifier) > 63 {
		return "", fmt.Errorf("generated identifier '%s' with prefixSize is too long, this should never happen", identifier)
	}

	return identifier, nil
}

func dispatchMaxSizeByParameter(maxSize int, parameters ...string) map[string]int {
	parameterSizeMap := make(map[string]int)
	parameterBaseSize := maxSize / len(parameters)
	rest := maxSize % len(parameters)

	// Add to the rest when the parameter is shorter than the parameter max size.
	for _, parameter := range parameters {
		if parameterBaseSize > len(parameter) {
			rest += parameterBaseSize - len(parameter)
		}
	}

	// Dispatch the rest over the parameters in slice order
	for _, parameter := range parameters {
		size := parameterBaseSize + rest

		if size > len(parameter) {
			// remove the rest that is used for this parameter
			if parameterBaseSize <= len(parameter) {
				rest = size - len(parameter)
			}
			size = len(parameter)
		} else {
			rest = 0
		}
		parameterSizeMap[parameter] = size
	}
	return parameterSizeMap
}

func getStringBufferWithHeader(buf *bytes.Buffer, fileName string) (string, int) {
	out := buf.String()
	return fmt.Sprintf("-- File: %s\n%s", fileName, out), len(out)
}

// VisitFile prepare a .psql from a proto one
// For each messages, call VisitMessage
func (v *PSQLVisitor) VisitFile(f pgs.File) (pgs.Visitor, error) {
	v.Debugf("pssql: Processing file %s", f.Name().String())

	initializations := make([]string, 0)
	finalizations := make([]string, 0)

	ok, err := f.Extension(psql.E_Initialization, &initializations)
	if err != nil {
		v.Logf("Error can't retrieve initialization extensions for file %s with error: %s", f.Name().String(), err)
	} else if ok {
		_, err = v.initW.Write([]byte(strings.Join(initializations, "\n")))

		if err != nil {
			v.Logf("Error can't write initialization for file %s with error: %s", f.Name().String(), err)
		}
	}

	ok, err = f.Extension(psql.E_Finalization, &finalizations)
	if err != nil {
		v.Logf("Error can't retrieve finalization extensions for file %s with error: %s", f.Name().String(), err)
	} else if ok {
		_, err = v.initW.Write([]byte(strings.Join(initializations, "\n")))

		if err != nil {
			v.Logf("Error can't write finalization for file %s with error: %s", f.Name().String(), err)
		}
	}

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
	relayCascadeUpdates := make([]*psql.RelayCascadeUpdate, 0)

	if ok, err := m.Extension(psql.E_Prefix, &prefixes); ok && err != nil {
		v.Logf("Error can't retrieve prefix extensions for message %s with error: %s", m.Name().String(), err)
	}

	if ok, err := m.Extension(psql.E_Constraint, &constraints); ok && err != nil {
		v.Logf("Error can't retrieve constraint extensions for message %s with error: %s", m.Name().String(), err)
	}

	if ok, err := m.Extension(psql.E_Suffix, &suffixes); ok && err != nil {
		v.Logf("Error can't retrieve suffix extensions for message %s with error: %s", m.Name().String(), err)
	}

	ok, err = m.Extension(psql.E_RelayCascadeUpdate, &relayCascadeUpdates)
	if err != nil {
		v.Logf("Error can't retrieve relay cascades updates options for message %s with error: %s", m.Name().String(), err)
	} else if ok {
		v.writeRelayCascadeUpdate(m.Name().String(), relayCascadeUpdates)
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

	definitions := appendSlices(prefixes, v.columns, constraints, suffixes)

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
	var cascadeUpdateOnRelatedTables []*psql.CascadeUpdateOnRelatedTable

	ok, err := f.Extension(psql.E_Column, &column)
	if ok && err == nil {
		v.columns = append(v.columns, f.Name().String()+" "+column)
	}

	ok, err = f.Extension(psql.E_AutoFillOnUpdate, &auto_fill_value)
	if ok && err == nil {
		v.writeAutoFillUpdate(f.Message().Name().String(), f.Name().String(), auto_fill_value)
	}

	ok, err = f.Extension(psql.E_CascadeUpdateOnRelatedTable, &cascadeUpdateOnRelatedTables)
	if ok && err == nil {
		v.writeCascadeUpdateOnRelatedTable(f.Message().Name().String(), f.Name().String(), cascadeUpdateOnRelatedTables)
	}

	return nil, nil
}
