package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_createLetterSlices(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want map[int][]string
	}{
		{
			name: "Split into 4 batches",
			args: args{
				n: 4,
			},
			want: map[int][]string{
				0: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m"},
				1: []string{"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
				2: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M"},
				3: []string{"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"},
			},
		},
		{
			name: "Split into 5 batches",
			args: args{
				n: 5,
			},
			want: map[int][]string{
				0: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
				1: []string{"k", "l", "m", "n", "o", "p", "q", "r", "s", "t"},
				2: []string{"u", "v", "w", "x", "y", "z", "A", "B", "C", "D"},
				3: []string{"E", "F", "G", "H", "I", "J", "K", "L", "M", "N"},
				4: []string{"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X"},
				5: []string{"Y", "Z"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createLetterSlices(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createLetterSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_separateWords(t *testing.T) {
	type args struct {
		n     int
		esh   string
		lines []string
	}
	tests := []struct {
		name string
		args args
		want map[int][]string
	}{
		{
			name: "Processing 2 words without last encrypted empty line",
			args: args{
				n:     4,
				esh:   "AAAAA",
				lines: []string{"H", "He", "Hel", "Hell", "Hello", "AAAAA", "L", "Lu", "Luc", "Luci", "Lucie", "Lucien"},
			},
			want: map[int][]string{
				0: []string{"H", "He", "Hel", "Hell", "Hello"},
				1: []string{"L", "Lu", "Luc", "Luci", "Lucie", "Lucien"},
			},
		},
		{
			name: "Processing 2 words with last encrypted empty line",
			args: args{
				n:     4,
				esh:   "AAAAA",
				lines: []string{"H", "He", "Hel", "Hell", "Hello", "AAAAA", "L", "Lu", "Luc", "Luci", "Lucie", "Lucien", "AAAAA"},
			},
			want: map[int][]string{
				0: []string{"H", "He", "Hel", "Hell", "Hello"},
				1: []string{"L", "Lu", "Luc", "Luci", "Lucie", "Lucien"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := separateWords(tt.args.n, tt.args.esh, tt.args.lines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("separateWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decipherLine(t *testing.T) {
	type args struct {
		done        chan struct{}
		decryptedCh chan string
		concCount   int
		encryptF    encryptFunc
		letters     []string
		decrypted   string
		encrypted   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				done:        make(chan struct{}, 2),
				decryptedCh: make(chan string, 1),
				concCount:   2,
				encryptF:    encryptSHA1(),
				letters:     []string{"a", "b", "c", "d", "e", "f"},
				decrypted:   "c",
				encrypted:   "b452d6b23b3c28f85872fffd99bdaf90ce0ad44a",
			},
			want: "ce",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decipherLine(tt.args.done, tt.args.decryptedCh, tt.args.concCount, tt.args.encryptF, tt.args.letters, tt.args.decrypted, tt.args.encrypted)

			select {
			case <-time.After(1 * time.Second):
				t.Errorf("Timeout in %s", tt.name)
			case got := <-tt.args.decryptedCh:
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("separateWords() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
