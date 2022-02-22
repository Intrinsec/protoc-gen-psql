package main

import (
	_ "embed"
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
