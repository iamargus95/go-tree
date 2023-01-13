package main

import (
	"bytes"
	"testing"
)

func TestOptions_tree(t *testing.T) {
	type fields struct {
		Indent    string
		OutputBuf string
		ShowFiles bool
	}
	type args struct {
		path   string
		isLast bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			args: args{
				path:   "./test/success",
				isLast: false,
			},
			fields: fields{
				Indent:    "",
				OutputBuf: "./test/success\n├── 1.txt\n├── 2.txt\n└── 3.txt\n0 directories , 3 files",
				ShowFiles: false,
			},
		},
		{
			name: "nested empty dir",
			args: args{
				path:   "./test/empty",
				isLast: false,
			},
			fields: fields{
				Indent:    "",
				OutputBuf: "./test/empty\n└── empty2\n1 directories , 0 files",
				ShowFiles: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testBuf := &bytes.Buffer{}
			testBuf.WriteString(tt.fields.OutputBuf)

			o := &Options{
				Indent:    tt.fields.Indent,
				OutputBuf: &bytes.Buffer{},
				ShowFiles: tt.fields.ShowFiles,
			}
			o.printTree(tt.args.path)

			if testBuf.String() != o.OutputBuf.String() {
				t.Fatalf("GOT = \n%s \nWANT = \n%s", o.OutputBuf.String(), testBuf.String())
			}
		})
	}
}
