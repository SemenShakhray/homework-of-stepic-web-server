package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"
)

type ParamsMethod struct {
	URL        string `json:"url"`
	Auth       bool   `json:"auth"`
	Method     string `json:"method"`
	Reciever   string
	NameMethod string
}

type FieldTagStruct struct {
	NameField string
	Tag       string
}

type StructTags struct {
	NameStruct string
	TagsFields []FieldTagStruct
}

type Model struct {
	Field string
	Tag   string
}

type ModelResult struct {
	Name       string
	FieldsTags []Model
}

type Collection struct {
	Params []ParamsMethod
	Tags   []StructTags
	Model  []ModelResult
}

type tmplParam struct {
	NameFieldLow string
	Min          string
	Max          string
	Required     bool
	ParamName    string
	Enum         string
	Default      string
	NameField    string
}

var tmplServeHTTP = `func (h *{{.Reciever}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path { {{range .Res}}
	case "{{.URL}}": 
	h.{{.NameMethod}}User(w, r) {{end}}
	default:
		http.Error(w, "{\"error\":\"unknown method\"}", 404)
	}
}
	`
var tmplValidParam = `
func Valid{{.StructName}}(r *http.Request) ({{.StructName}}, error) {
	in:={{.StructName}}{} {{range .Param}} 
	{{if .Max}}
	{{.NameField}}, err:=strconv.Atoi(r.FormValue("{{.NameFieldLow}}"))
	if err != nil {
			return in, fmt.Errorf("{\"error\":\"{{.NameFieldLow}} must be int\"}")
		} 
			in.{{.NameField}}={{.NameField}} {{else if .ParamName}}
	in.{{.NameField}} = r.FormValue("{{.ParamName}}") {{else}}
	in.{{.NameField}} = r.FormValue("{{.NameFieldLow}}") {{end}} {{if .Required}}
		if in.{{.NameField}} == "" {
			return in, fmt.Errorf("{\"error\":\"{{.NameFieldLow}} must me not empty\"}")
		} {{end}} {{if and .Min (not .Max)}}
		if len(in.{{.NameField}}) < {{.Min}} {
			return in, fmt.Errorf("{\"error\":\"{{.NameFieldLow}} len must be >= {{.Min}}\"}")
		} {{end}} {{if and .Min .Max}}
		 if in.{{.NameField}} < {{.Min}} {
			return in, fmt.Errorf("{\"error\":\"{{.NameFieldLow}} must be >= {{.Min}}\"}")
	}
			if in.{{.NameField}}>{{.Max}} {
	return in, fmt.Errorf("{\"error\":\"{{.NameFieldLow}} must be <= {{.Max}}\"}")		
	} {{end}} {{if .Enum}}
	if in.{{.NameField}}=="" {
	in.{{.NameField}}="{{.Default}}"
	}
	enum:=strings.Split("{{.Enum}}", "|")
	mapEnum:=make(map[string]string)
	for _,v:=range enum {
	mapEnum[v]=""
	}
	if _, ok:=mapEnum[in.{{.NameField}}]; !ok {
	return in, fmt.Errorf("{\"error\":\"{{.NameFieldLow}} must be one of [%s, %s, %s]\"}",enum[0], enum[1], enum[2])		
	} {{end}} {{end}}
	return in, nil
	}
	`

var tmplWrapper = `
func (h *{{.Reciever}}) {{.NameMethod}}User(w http.ResponseWriter, r *http.Request) {
		{{if eq .Method "POST"}}
		if r.Method == http.MethodPost {
			password := r.Header.Get("X-Auth")
		if password == "" {
			http.Error(w, "{\"error\":\"unauthorized\"}", 403)
			return
		} 
				if password != "100500" {
			http.Error(w, "{\"error\":\"wrong password\"}", http.StatusUnauthorized)
			return
		}{{end}}
		ctx := context.TODO()
		in, err:=Valid{{.StructName}}(r)
		if err != nil {
		http.Error(w, err.Error(), 400)
		return
		}
		res, err := h.{{.NameMethod}}(ctx, in)
		if err != nil {
			ae, ok := err.(ApiError)
			if !ok {
				http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), http.StatusInternalServerError)
				return
			}
			http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err.Error()), ae.HTTPStatus)
			return
		}
		result := Result{
			"error": "",
			"response": Result{
			{{range .Model}} 
			"{{.Tag}}": res.{{.Field}}, {{end}}
			},
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, "error encoding JSON", http.StatusInternalServerError)
		}
	 	{{if eq .Method "POST"}} } else {
			http.Error(w, "{\"error\":\"bad method\"}", 406)
			} {{end}}
	}`

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	)`)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `type Result map[string]interface{}`)
	fmt.Fprintln(out) // empty line
	var result = &Collection{}
	dataCollection(node, result)
	addServeHTTP(out, result)
	CreateValidator(out, result)
	CreateWrappers(out, result)
}

func dataCollection(node *ast.File, result *Collection) {
	nameParamMethod := make(map[string]string)
	nameModel := make(map[string]string)
	for _, field := range node.Decls {
		doc := &ParamsMethod{}
		f, ok := field.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if f.Doc == nil {
			continue
		}
		comment := strings.TrimPrefix(f.Doc.Text(), "apigen:api ")
		err := json.Unmarshal([]byte(comment), doc)
		if err != nil {
			fmt.Println(err)
		}
		doc.Reciever = f.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
		doc.NameMethod = fmt.Sprint(f.Name)
		for _, item := range f.Type.Params.List {
			nameType, ok := item.Type.(*ast.Ident)
			if !ok {
				continue
			}
			nameParamMethod[fmt.Sprint(nameType)] = ""
		}
		for _, item := range f.Type.Results.List {
			typeStar, ok := item.Type.(*ast.StarExpr)
			if !ok {
				continue
			}
			nameType, ok := typeStar.X.(*ast.Ident)
			if !ok {
				continue
			}
			nameModel[fmt.Sprint(nameType)] = ""
		}
		result.Params = append(result.Params, *doc)
	}

	for _, field := range node.Decls {
		tags := &StructTags{}
		str, ok := field.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range str.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structName := currType.Name.Name

			if _, ok := nameParamMethod[structName]; !ok {
				continue
			}
			tags.NameStruct = structName
			for _, t := range currType.Type.(*ast.StructType).Fields.List {
				tf := FieldTagStruct{
					NameField: t.Names[0].Name,
					Tag:       t.Tag.Value,
				}
				tags.TagsFields = append(tags.TagsFields, tf)
			}
			result.Tags = append(result.Tags, *tags)
		}
	}

	for _, field := range node.Decls {
		tags := &ModelResult{}
		str, ok := field.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range str.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structName := currType.Name.Name

			if _, ok := nameModel[structName]; !ok {
				continue
			}
			tags.Name = structName
			for _, t := range currType.Type.(*ast.StructType).Fields.List {
				f := reflect.StructTag(strings.Trim(t.Tag.Value, "`"))
				tag := f.Get("json")
				tf := Model{
					Field: t.Names[0].Name,
					Tag:   tag,
				}
				tags.FieldsTags = append(tags.FieldsTags, tf)
			}
			result.Model = append(result.Model, *tags)
		}
	}
}

func addServeHTTP(out io.Writer, result *Collection) {
	type ParamsSearveHTTP struct {
		Reciever string
		Res      []ParamsMethod
	}
	t := template.New("example")
	t, _ = t.Parse(tmplServeHTTP)
	var Recievers []string
	mapRec := make(map[string]string)
	for _, item := range result.Params {
		if _, ok := mapRec[item.Reciever]; !ok {
			Recievers = append(Recievers, item.Reciever)
			mapRec[item.Reciever] = ""
		}
	}
	for _, rec := range Recievers {
		param := ParamsSearveHTTP{}
		param.Reciever = rec
		for _, item := range result.Params {
			if rec == item.Reciever {
				param.Res = append(param.Res, item)
			} else {
				continue
			}
		}
		t.Execute(out, param)
	}

}

func checkTag(values []string) tmplParam {
	tmplP := tmplParam{}
	for _, v := range values {
		if strings.Contains(v, "required") {
			tmplP.Required = true
		}
		if strings.Contains(v, "paramname") {
			p := strings.Split(v, "=")
			tmplP.ParamName = p[1]
		}
		if strings.Contains(v, "min") && !strings.Contains(v, "admin") {
			p := strings.Split(v, "=")
			tmplP.Min = p[1]
		}
		if strings.Contains(v, "max") {
			p := strings.Split(v, "=")
			tmplP.Max = p[1]
		}
		if strings.Contains(v, "enum") {
			s := strings.TrimPrefix(v, "enum=")
			tmplP.Enum = s
		}
		if strings.Contains(v, "default") {
			p := strings.Split(v, "=")
			tmplP.Default = p[1]
		}
	}
	return tmplP
}

func CreateValidator(out io.Writer, result *Collection) {
	type AllParam struct {
		StructName string
		Param      []tmplParam
	}
	tmpl := template.New("validation")
	tmpl.Parse(tmplValidParam)
	for _, item := range result.Tags {
		p := AllParam{}
		for _, t := range item.TagsFields {
			f := reflect.StructTag(strings.Trim(t.Tag, "`"))
			tag := f.Get("apivalidator")
			values := strings.Split(tag, ",")
			paramOfTmpl := checkTag(values)
			paramOfTmpl.NameField = t.NameField
			paramOfTmpl.NameFieldLow = strings.ToLower(t.NameField)
			p.StructName = item.NameStruct
			p.Param = append(p.Param, paramOfTmpl)
		}
		tmpl.Execute(out, p)
	}
}

func CreateWrappers(out io.Writer, result *Collection) {
	type param struct {
		StructName string
		Reciever   string
		Method     string
		NameMethod string
		Model      []Model
	}

	tmpl, err := template.New("wrapper").Parse(tmplWrapper)
	if err != nil {
		log.Fatal(err)
	}

	var nameStruct []string
	for _, item := range result.Tags {
		nameStruct = append(nameStruct, item.NameStruct)
	}
	for i, item := range result.Params {
		p := param{}
		p.StructName = nameStruct[i]
		p.Method = item.Method
		p.NameMethod = item.NameMethod
		p.Reciever = item.Reciever
		p.Model = result.Model[i].FieldsTags
		fmt.Println(p)
		tmpl.Execute(out, p)
	}
}

// находясь в папке выше
// go build -o ./codegen.exe gen/* && ./codegen.exe pack/unkack.go  pack/marshaller.go
// go run pack/*
