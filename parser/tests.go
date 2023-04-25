package parser

import (
	"bytes"
	"html/template"
	"strings"
)

type Tests struct {
	list     []*Test
	funcName string
	params   []*param
}

func (_this *Tests) FuncName() string {
	return _this.funcName
}

func newTests(list []*Test, funcName string, params []*param) *Tests {
	return &Tests{list, funcName, params}
}

const codeTmpl = `import (
	"reflect"
	"testing"
)

func Test_{{.FuncName}}(t *testing.T) {
	type args struct {
		{{- range .Input}}
		{{.}}{{end}}
	}
	tests := []struct {
		name string
		args args
		{{if gt (len .Output) 0 }}{{"want"}} {{.Output}}{{end}}
	}{
{{.Tests}}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := {{.FuncName}}({{.FuncParams}}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("{{.FuncName}}() = %v, want %v", got, tt.want)
			}
		})
	}
}`

func (_this *Tests) buffer() *bytes.Buffer {
	output := ""
	var input, fp []string
	for _, p := range _this.params {
		if p.key == "output" {
			output = p.val
			continue
		}
		input = append(input, p.key+" "+p.val)
		fp = append(fp, "tt.args."+p.key)
	}

	tests := strings.Builder{}
	end := len(_this.list) - 1
	for i, test := range _this.list {
		var args, want = []byte("args{"), []byte{}
		last := len(test.Value()) - 1
		for j, code := range test.Value() {
			if j == last && output != "" {
				want = code.Value()
				continue
			}
			args = append(args, code.Value()...)
			args = append(args, ',')
		}
		args[len(args)-1] = '}'
		args = append(args, ',', ' ')
		name := append([]byte("{\"test"), byte(i+'1'), '"', ',', ' ')
		want = append(want, '}', ',')
		tests.WriteByte('\t')
		tests.WriteByte('\t')
		tests.Write(name)
		tests.Write(args)
		tests.Write(want)
		if i == end {
			continue
		}
		tests.WriteByte('\n')
	}

	t, _ := template.New("").Parse(codeTmpl)
	w := &bytes.Buffer{}
	_ = t.Execute(w, struct {
		Input      []string
		Output     string
		Tests      template.HTML
		FuncName   string
		FuncParams string
	}{
		input, output, template.HTML(tests.String()), _this.funcName, strings.Join(fp, ", "),
	})
	return w
}

func (_this *Tests) ToCode() string {
	return _this.buffer().String()
}

func (_this *Tests) ToBytes() []byte {
	return _this.buffer().Bytes()
}
