package main

import "testing"

func Test_decryptLines(t *testing.T) {
	type args struct {
		emptyEncr string
		decrypt   map[string]string
		lines     []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "One word ending with hashed empty string",
			args: args{
				emptyEncr: "zzz",
				decrypt: map[string]string{
					"AAA": "a",
					"BBB": "b",
					"CCC": "c",
					"DDD": "d",
					"EEE": "e",
				},
				lines: []string{"zzz", "EEE", "BBB", "AAA", "CCC", "zzz"},
			},
			want:    "ebac",
			wantErr: false,
		},
		{
			name: "One word not ending hashed empty string",
			args: args{
				emptyEncr: "zzz",
				decrypt: map[string]string{
					"AAA": "a",
					"BBB": "b",
					"CCC": "c",
					"DDD": "d",
					"EEE": "e",
				},
				lines: []string{"zzz", "EEE", "BBB", "AAA", "CCC"},
			},
			want:    "ebac",
			wantErr: false,
		},
		{
			name: "Two words ending with hashed empty string",
			args: args{
				emptyEncr: "zzz",
				decrypt: map[string]string{
					"AAA": "a",
					"BBB": "b",
					"CCC": "c",
					"DDD": "d",
					"EEE": "e",
				},
				lines: []string{"zzz", "EEE", "BBB", "AAA", "CCC", "zzz", "DDD", "AAA", "BBB", "zzz"},
			},
			want:    "ebac dab",
			wantErr: false,
		},
		{
			name: "Two words not ending with hashed empty string",
			args: args{
				emptyEncr: "zzz",
				decrypt: map[string]string{
					"AAA": "a",
					"BBB": "b",
					"CCC": "c",
					"DDD": "d",
					"EEE": "e",
				},
				lines: []string{"zzz", "EEE", "BBB", "AAA", "CCC", "zzz", "DDD", "AAA", "BBB"},
			},
			want:    "ebac dab",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decryptLines(tt.args.emptyEncr, tt.args.decrypt, tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("decryptLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decryptLines() = [%v], want [%v]", got, tt.want)
			}
		})
	}
}
