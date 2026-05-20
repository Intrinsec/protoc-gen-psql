package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func Test_generateIdentifierName(t *testing.T) {
	type args struct {
		name       string
		prefixSize int
		parameters []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "2 long parameters",
			args: args{
				name:       "cascade_related",
				prefixSize: 10,
				parameters: []string{
					"Communication",
					"raw_communication_uuid",
				},
			},
			want: "cascade_related_communication_raw_communicat_d500eb9",
		},
		{
			name: "3 parameters with first one below max size",
			args: args{
				name:       "cascade_related",
				prefixSize: 10,
				parameters: []string{
					"action",
					"communication_uuid",
					"incident",
				},
			},
			want: "cascade_related_action_communicatio_incident_ef530db3",
		},
		{
			name: "3 parameters with first one above max size",
			args: args{
				name:       "cascade_related",
				prefixSize: 10,
				parameters: []string{
					"action_id",
					"communication_uuid",
					"incident",
				},
			},
			want: "cascade_related_action_id_communica_incident_1a020edf",
		},
		{
			name: "3 parameters, all below or equal to max size",
			args: args{
				name:       "cascade_related",
				prefixSize: 10,
				parameters: []string{
					"id",
					"inc_uuid",
					"com_uuid",
				},
			},
			want: "cascade_related_id_inc_uuid_com_uuid_501607cd",
		},
		{
			name: "prefix size too large returns error",
			args: args{
				name:       "x",
				prefixSize: 99,
				parameters: []string{"a"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateIdentifierName(tt.args.name, tt.args.prefixSize, tt.args.parameters...)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateIdentifierName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateIdentifierName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateCascadeIdentifierNames(t *testing.T) {
	fn, tg, tgDel, fnCreate := generateCascadeIdentifierNames("auto_fill", "Incident", "update_time")

	for _, got := range []string{fn, tg, tgDel, fnCreate} {
		if len(got) > 63 {
			t.Errorf("identifier %q exceeds 63 characters", got)
		}
	}

	cases := map[string]string{
		"fn":       fn,
		"tg":       tg,
		"tgDel":    tgDel,
		"fnCreate": fnCreate,
	}
	expectedPrefixes := map[string]string{
		"fn":       "fn_",
		"tg":       "zz_tg_",
		"tgDel":    "tg_del_",
		"fnCreate": "fn_create_",
	}
	for key, ident := range cases {
		if !strings.HasPrefix(ident, expectedPrefixes[key]) {
			t.Errorf("%s identifier %q missing prefix %q", key, ident, expectedPrefixes[key])
		}
	}
}

func Test_allocateRoomToParameters(t *testing.T) {
	tests := []struct {
		name       string
		maxSize    int
		parameters []string
		check      func(t *testing.T, got map[string]int)
	}{
		{
			name:       "equal split when all parameters longer than slot",
			maxSize:    20,
			parameters: []string{"abcdefghij", "klmnopqrst"},
			check: func(t *testing.T, got map[string]int) {
				total := 0
				for _, n := range got {
					total += n
				}
				if total != 20 {
					t.Errorf("total allocated = %d, want 20", total)
				}
			},
		},
		{
			name:       "short parameters do not exceed their own length",
			maxSize:    30,
			parameters: []string{"a", "bb", "ccc"},
			check: func(t *testing.T, got map[string]int) {
				lengths := map[string]int{"a": 1, "bb": 2, "ccc": 3}
				for p, n := range got {
					if n > lengths[p] {
						t.Errorf("parameter %q got %d slots but is only %d chars", p, n, lengths[p])
					}
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allocateRoomToParameters(tt.maxSize, tt.parameters...)
			tt.check(t, got)
		})
	}
}

func Test_appendSlices(t *testing.T) {
	tests := []struct {
		name string
		in   [][]string
		want []string
	}{
		{
			name: "two non-empty slices",
			in:   [][]string{{"a", "b"}, {"c"}},
			want: []string{"a", "b", "c"},
		},
		{
			name: "empty inputs return empty result",
			in:   [][]string{nil, {}, nil},
			want: nil,
		},
		{
			name: "single slice passes through",
			in:   [][]string{{"x", "y", "z"}},
			want: []string{"x", "y", "z"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendSlices(tt.in...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringBufferWithHeader(t *testing.T) {
	buf := bytes.NewBufferString("CREATE TABLE x();")
	out, count := getStringBufferWithHeader(buf, "x.proto")

	if count != len("CREATE TABLE x();") {
		t.Errorf("count = %d, want %d", count, len("CREATE TABLE x();"))
	}
	if !strings.HasPrefix(out, "-- File: x.proto\n") {
		t.Errorf("output missing header prefix; got %q", out)
	}
	if !strings.HasSuffix(out, "CREATE TABLE x();") {
		t.Errorf("output missing original buffer content; got %q", out)
	}
}

func Test_getStringBufferWithHeader_emptyBuffer(t *testing.T) {
	buf := &bytes.Buffer{}
	out, count := getStringBufferWithHeader(buf, "empty.proto")

	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}
	if out != "-- File: empty.proto\n" {
		t.Errorf("unexpected header-only output: %q", out)
	}
}

func Test_generateFromTemplate(t *testing.T) {
	tmpl := "name={{.Name}};join={{JoinedStrings .Items \"|\"}}"
	data := struct {
		Name  string
		Items []string
	}{
		Name:  "test",
		Items: []string{"a", "b", "c"},
	}
	var buf bytes.Buffer
	generateFromTemplate(tmpl, "fixture", data, &buf)
	got := buf.String()
	want := "name=test;join=a|b|c"
	if got != want {
		t.Errorf("generateFromTemplate() rendered %q, want %q", got, want)
	}
}
