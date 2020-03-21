package py

import (
	"bytes"
	"fmt"
	"github.com/go-openapi/spec"
	"github.com/iancoleman/strcase"
	"strings"
)

func generateDefinitions(pkg string, sg *spec.Swagger) (p []byte, err error) {

	bb := bytes.NewBufferString("")
	for name, def := range sg.Definitions {

		if name == "commons.ErrorResult" {
			continue
		}

		generateEntry(pkg, bb, name, def)

		generateEntryList(pkg, bb, name, def)

	}

	p = bb.Bytes()

	return
}

type Prop struct {
	Name      string
	Kind      string
	ClassName string
}

func (p *Prop) String() string {
	return p.Name
}

func generateEntry(pkg string, bb *bytes.Buffer, name string, def spec.Schema)  {
	className := name[strings.LastIndex(name, ".")+1:]
	className = strcase.ToCamel(className)

	arrayKind := ""
	if def.Type.Contains("array") {
		arrayKind = "List"
	}

	bb.WriteString(fmt.Sprintf("class %s%s(object):\n", className, arrayKind))

	// description
	bb.WriteString(`    """`)
	bb.WriteString("\n")
	bb.WriteString(fmt.Sprintf(`    %s`, def.Description))
	bb.WriteString("\n")
	bb.WriteString(`    """`)
	bb.WriteString("\n")

	// Properties
	props := def.Properties

	if props == nil || len(props) == 0 {
		return
	}

	props0 := make([]string, 0, 1)
	propKeys := make([]*Prop, 0, 1)

	for propName, prop := range props {
		propName = strcase.ToSnake(propName)
		if prop.Type.Contains("object") {

			ref := prop.Ref.String()
			ref = strings.ReplaceAll(ref, "/", ".")
			refName := ref[strings.LastIndex(ref, ".")+1:]
			refName = strcase.ToCamel(refName)

			props00 := fmt.Sprintf("%s: %s.%s = None", propName, pkg, refName)
			bb.WriteString(fmt.Sprintf("    %s\n", props00))

			props0 = append(props0, props00)
			propKeys = append(propKeys, &Prop{Name: propName, Kind: "object", ClassName: refName})
			continue
		}
		if prop.Type.Contains("array") {

			ref := prop.Items.Schema.Ref.String()
			ref = strings.ReplaceAll(ref, "/", ".")
			refName := ref[strings.LastIndex(ref, ".")+1:]
			refName = strcase.ToCamel(refName)

			props00 := fmt.Sprintf("%s: %s.%sList = None", propName, pkg, refName)
			bb.WriteString(fmt.Sprintf("    %s\n", props00))

			props0 = append(props0, props00)
			propKeys = append(propKeys, &Prop{Name: propName, Kind: "array", ClassName: fmt.Sprintf("%sList", refName)})
			continue
		}

		if prop.Type.Contains("integer") {
			props00 := fmt.Sprintf("%s: int = None", propName)
			bb.WriteString(fmt.Sprintf("    %s\n", props00))
			props0 = append(props0, props00)
			propKeys = append(propKeys, &Prop{Name: propName, Kind: "int"})
			continue
		}

		if prop.Type.Contains("number") {
			props00 := fmt.Sprintf("%s: float = None", propName)
			bb.WriteString(fmt.Sprintf("    %s\n", props00))
			props0 = append(props0, props00)
			propKeys = append(propKeys, &Prop{Name: propName, Kind: "float"})
			continue
		}

		if prop.Type.Contains("boolean") {
			props00 := fmt.Sprintf("%s: bool = None", propName)
			bb.WriteString(fmt.Sprintf("    %s\n", props00))
			props0 = append(props0, props00)
			propKeys = append(propKeys, &Prop{Name: propName, Kind: "bool"})
			continue
		}

		if prop.Type.Contains("string") {
			format := strings.TrimSpace(prop.Format)
			if format == "byte" {
				props00 := fmt.Sprintf("%s: chr = None", propName)
				bb.WriteString(fmt.Sprintf("    %s\n", props00))
				props0 = append(props0, props00)
				propKeys = append(propKeys, &Prop{Name: propName, Kind: "chr"})
				continue
			}
			if format == "binary" {
				props00 := fmt.Sprintf("%s: bytes = None", propName)
				bb.WriteString(fmt.Sprintf("    %s\n", props00))
				props0 = append(props0, props00)
				propKeys = append(propKeys, &Prop{Name: propName, Kind: "bytes"})
				continue
			}
			if format == "date-time" {
				props00 := fmt.Sprintf("%s: datetime.datetime = None", propName)
				bb.WriteString(fmt.Sprintf("    %s\n", props00))
				props0 = append(props0, props00)
				propKeys = append(propKeys, &Prop{Name: propName, Kind: "datetime.datetime"})
				continue
			}
			props00 := fmt.Sprintf("%s: str = None", propName)
			bb.WriteString(fmt.Sprintf("    %s\n", props00))
			props0 = append(props0, props00)
			propKeys = append(propKeys, &Prop{Name: propName, Kind: "str"})
			continue
		}

	}

	bb.WriteString("\n")

	// __init__

	bb.WriteString(`    def __init__(self`)
	initArgs := ""
	if len(props0) > 0 {
		bbInitArgs := bytes.NewBufferString("")
		for _, props00 := range props0 {
			bbInitArgs.WriteString(", ")
			bbInitArgs.WriteString(props00)
		}
		initArgs = bbInitArgs.String()

	}
	bb.WriteString(initArgs)
	bb.WriteString("):")
	bb.WriteString("\n")

	if len(propKeys) > 0 {
		for _, propsKey := range propKeys {
			bb.WriteString(fmt.Sprintf("        self.%s = %s\n", propsKey.Name, propsKey.Name))
		}
	}
	bb.WriteString("\n")

	if len(propKeys) > 0 {
		// from_dict
		bb.WriteString(`    def from_dict(self, data: dict):`)
		bb.WriteString("\n")

		for _, propsKey := range propKeys {
			bb.WriteString(fmt.Sprintf(`        self.%s = data["%s"]`, propsKey.Name, propsKey.Name))
			bb.WriteString("\n")
		}
		bb.WriteString(`        return self`)
		bb.WriteString("\n")
		bb.WriteString("\n")

		// to_dict
		bb.WriteString(`    def to_dict(self):`)
		bb.WriteString("\n")
		dictBytes := bytes.NewBufferString("")
		for _, propsKey := range propKeys {
			dictBytes.WriteString(fmt.Sprintf(`, "%s": self.%s`, propsKey.Name, propsKey.Name))
		}
		bb.WriteString(fmt.Sprintf(`        return {%s}`, dictBytes.String()[2:]))
		bb.WriteString("\n")
		bb.WriteString("\n")

		// from_json
		bb.WriteString(`    def from_json(self, data: json):`)
		bb.WriteString("\n")
		bb.WriteString("        for k, v in data.items():")
		bb.WriteString("\n")
		for i, propsKey := range propKeys {
			if i == 0 {
				bb.WriteString(fmt.Sprintf(`            if k == "%s":`, propsKey.Name))
				bb.WriteString("\n")
				if propsKey.Kind == "array" {
					bb.WriteString(fmt.Sprintf(`                self.%s = %s.%s().from_json_array(v)`, propsKey.Name, pkg, propsKey.ClassName))
				} else if propsKey.Kind == "object" {
					bb.WriteString(fmt.Sprintf(`                self.%s = %s.%s().from_json_object(v)`, propsKey.Name,pkg, propsKey.ClassName))
				} else {
					bb.WriteString(fmt.Sprintf(`                self.%s = v`, propsKey))
				}
				bb.WriteString("\n")
				continue
			}
			bb.WriteString(fmt.Sprintf(`            elif k == "%s":`, propsKey.Name))
			bb.WriteString("\n")
			if propsKey.Kind == "array" {
				bb.WriteString(fmt.Sprintf(`                self.%s = %s.%s().from_json_array(v)`, propsKey.Name,pkg, propsKey.ClassName))
			} else if propsKey.Kind == "object" {
				bb.WriteString(fmt.Sprintf(`                self.%s = %s.%s().from_json_object(v)`, propsKey.Name,pkg, propsKey.ClassName))
			} else {
				bb.WriteString(fmt.Sprintf(`                self.%s = v`, propsKey))
			}
			bb.WriteString("\n")
			continue

		}
		bb.WriteString(`        return self`)
		bb.WriteString("\n")
		bb.WriteString("\n")

		// to_json
		bb.WriteString(`    def to_json(self):`)
		bb.WriteString("\n")
		bb.WriteString("        return json.dumps(self.to_dict())")
		bb.WriteString("\n")
		bb.WriteString("\n")

		// __str__
		bb.WriteString(`    def __str__(self):`)
		bb.WriteString("\n")
		bb.WriteString("        return self.to_dict().__str__()")
		bb.WriteString("\n")
		bb.WriteString("\n")
	}

	bb.WriteString("\n")
}

func generateEntryList(pkg string, bb *bytes.Buffer, name string, def spec.Schema) {

	className := name[strings.LastIndex(name, ".")+1:]
	className = strcase.ToCamel(className)

	bb.WriteString(fmt.Sprintf("class %sList(object):\n", className))

	// description
	bb.WriteString(`    """`)
	bb.WriteString("\n")
	bb.WriteString(fmt.Sprintf(`    %s`, def.Description))
	bb.WriteString("\n")
	bb.WriteString(`    """`)
	bb.WriteString("\n")

	bb.WriteString(fmt.Sprintf("    array: [%s.%s] = None", pkg, className))
	bb.WriteString("\n")
	bb.WriteString("\n")

	// init
	bb.WriteString(fmt.Sprintf("    def __init__(self, value: [%s.%s] = None):",pkg, className))
	bb.WriteString("\n")
	bb.WriteString("        if value is None:")
	bb.WriteString("\n")
	bb.WriteString("            self.array = None")
	bb.WriteString("\n")
	bb.WriteString("        else:")
	bb.WriteString("\n")
	bb.WriteString("            self.array = value")
	bb.WriteString("\n")
	bb.WriteString("\n")

	// from_dict_array
	bb.WriteString("    def from_dict_array(self, data: dict):")
	bb.WriteString("\n")
	bb.WriteString("        for d in data:")
	bb.WriteString("\n")
	bb.WriteString(fmt.Sprintf("            e = %s.%s()", pkg, className))
	bb.WriteString("\n")
	bb.WriteString("            e.from_dict(d)")
	bb.WriteString("\n")
	bb.WriteString("            self.array.append(e)")
	bb.WriteString("\n")
	bb.WriteString("        return self")
	bb.WriteString("\n")
	bb.WriteString("\n")

	// to_dict_array
	bb.WriteString("    def to_dict_array(self):")
	bb.WriteString("\n")
	bb.WriteString("        if self.array is None:")
	bb.WriteString("\n")
	bb.WriteString("            return None")
	bb.WriteString("\n")
	bb.WriteString("        xx = []")
	bb.WriteString("\n")
	bb.WriteString("        for x in self.array:")
	bb.WriteString("\n")
	bb.WriteString(fmt.Sprintf("            if isinstance(x, %s.%s):", pkg, className))
	bb.WriteString("\n")
	bb.WriteString("                xx.append(x.to_dict())")
	bb.WriteString("\n")
	bb.WriteString("        return xx")
	bb.WriteString("\n")
	bb.WriteString("\n")

	// from_json_array
	bb.WriteString("    def from_json_array(self, data: json):")
	bb.WriteString("\n")
	bb.WriteString("        if data is None:")
	bb.WriteString("\n")
	bb.WriteString("            return self")
	bb.WriteString("\n")
	bb.WriteString("        if len(data) == 0:")
	bb.WriteString("\n")
	bb.WriteString("            return self")
	bb.WriteString("\n")
	bb.WriteString("        self.array = []")
	bb.WriteString("\n")
	bb.WriteString("        _it = iter(data)")
	bb.WriteString("\n")
	bb.WriteString("        for _x in _it:")
	bb.WriteString("\n")
	bb.WriteString(fmt.Sprintf("            _e = %s.%s()", pkg, className))
	bb.WriteString("\n")
	bb.WriteString("            _e.from_dict(_x)")
	bb.WriteString("\n")
	bb.WriteString("            self.array.append(_e)")
	bb.WriteString("\n")
	bb.WriteString("        return self")
	bb.WriteString("\n")
	bb.WriteString("\n")

	// to_json_array
	bb.WriteString("    def to_json_array(self):")
	bb.WriteString("\n")
	bb.WriteString("        xx = []")
	bb.WriteString("\n")
	bb.WriteString("        for x in self.array:")
	bb.WriteString("\n")
	bb.WriteString("            xx.append(x.to_dict())")
	bb.WriteString("\n")
	bb.WriteString("        return json.dumps(xx)")
	bb.WriteString("\n")
	bb.WriteString("\n")

	// __str__
	bb.WriteString("    def __str__(self):")
	bb.WriteString("\n")
	bb.WriteString("        return self.to_dict_array().__str__()")
	bb.WriteString("\n")
	bb.WriteString("\n")

}