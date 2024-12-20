package utils

import (
	"fmt"
	"strings"
)

// Struct для входного JSON AtlasMapping
type AtlasMapping struct {
	AtlasMapping struct {
		DataSources []DataSource `json:"dataSource"`
		Mappings    Mappings     `json:"mappings"`
	} `json:"AtlasMapping"`
}

type DataSource struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Mappings struct {
	Mapping []Mapping `json:"mapping"`
}

type Mapping struct {
	ID          string      `json:"id"`
	InputFields []JsonField `json:"inputField,omitempty"`
	InputGroup  *FieldGroup `json:"inputFieldGroup,omitempty"`
	OutputField []JsonField `json:"outputField"`
	Priority    int         `json:"priority"`
}

type FieldGroup struct {
	Field []JsonField `json:"field"`
}

type Action struct {
	Index string `json:"index"`
	Type  string `json:"@type"`
}

type JsonField struct {
	Path      string   `json:"path"`
	Name      string   `json:"name"`
	FieldType string   `json:"fieldType"`
	DocID     string   `json:"docId"`
	Actions   []Action `json:"actions"`
}

// Структуры выходного JSON
type Output struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Steps       []Step `json:"steps"`
}

type Step struct {
	Name             string            `json:"name"`
	Mode             string            `json:"mode"`
	Val              Val               `json:"val"`
	Transform        bool              `json:"transform,omitempty"`
	RequestTransform *RequestTransform `json:"requestTransform,omitempty"`
}

type Val struct {
	Method string `json:"method,omitempty"`
	URL    string `json:"url"`
	Path   string `json:"path,omitempty"`
}

type RequestTransform struct {
	Spec map[string]string `json:"spec"`
}

// Генерация Val.url
func GetValURL(outputFields []JsonField, dataSourceMap map[string]string) string {
	for _, field := range outputFields {
		baseURL := dataSourceMap[field.DocID]
		return fmt.Sprintf("%s%s", baseURL, field.Path)
	}

	panic("DataSourceMap is not exist outputFields")
}

// Построение RequestTransform
func BuildRequestTransform(mapping Mapping) map[string]string {
	spec := make(map[string]string)
	if mapping.InputGroup != nil {
		for _, field := range mapping.InputGroup.Field {
			spec[strings.TrimLeft(field.Path, "/")] = fmt.Sprintf("%s%s", mapping.ID, strings.ReplaceAll(field.Path, "/", "."))
		}
	} else {
		for _, field := range mapping.InputFields {
			spec[strings.TrimLeft(field.Path, "/")] = fmt.Sprintf("%s%s", mapping.ID, strings.ReplaceAll(field.Path, "/", "."))
		}
	}
	return spec
}
