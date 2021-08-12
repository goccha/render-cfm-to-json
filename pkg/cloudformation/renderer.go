package cloudformation

import (
	"encoding/json"
	"github.com/awslabs/goformation/v5"
	"github.com/awslabs/goformation/v5/cloudformation"
	"github.com/awslabs/goformation/v5/intrinsics"
	"github.com/iancoleman/strcase"
	"golang.org/x/xerrors"
)

type Renderer struct {
	Name   string
	Params map[string]interface{}
}

func (ren *Renderer) Render(cfm string) (indent []byte, min []byte, err error) {
	var template *cloudformation.Template
	if template, err = openTemplate(cfm, ren.Params); err != nil {
		return nil, nil, err
	}
	var b []byte
	if b, err = template.JSON(); err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}
	body := make(map[string]interface{})
	if err = json.Unmarshal(b, &body); err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}
	body = toLowerCamelKey(body)
	if kv, ok := body["resources"]; ok {
		if resources, ok := kv.(map[string]interface{}); ok {
			indent, min, err = getProperties(ren.Name, resources)
		}
	}
	return
}

func openTemplate(cfm string, params map[string]interface{}) (template *cloudformation.Template, err error) {
	if params != nil {
		template, err = goformation.OpenWithOptions(cfm, &intrinsics.ProcessorOptions{
			ParameterOverrides: params,
		})
		if err != nil {
			err = xerrors.Errorf(": %w", err)
		}
		return
	}
	template, err = goformation.Open(cfm)
	if err != nil {
		err = xerrors.Errorf(": %w", err)
	}
	return
}

func toLowerCamelKey(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		k = strcase.ToLowerCamel(k)
		switch v := v.(type) {
		case map[string]interface{}:
			result[k] = toLowerCamelKey(v)
		case []interface{}:
			result[k] = toLowerCamelArray(v)
		default:
			result[k] = v
		}
	}
	return result
}

func toLowerCamelArray(array []interface{}) []interface{} {
	result := make([]interface{}, 0, len(array))
	for _, v := range array {
		switch v := v.(type) {
		case map[string]interface{}:
			result = append(result, toLowerCamelKey(v))
		case []interface{}:
			result = append(result, toLowerCamelArray(v))
		default:
			result = append(result, v)
		}
	}
	return result
}

func getProperties(name string, m map[string]interface{}) (indent []byte, min []byte, err error) {
	if body, ok := m[strcase.ToLowerCamel(name)]; ok {
		if m, ok := body.(map[string]interface{}); ok {
			if props, ok := m["properties"]; ok {
				if indent, err = json.MarshalIndent(props, "", "  "); err != nil {
					return nil, nil, xerrors.Errorf(": %w", err)
				}
				if min, err = json.Marshal(props); err != nil {
					return nil, nil, xerrors.Errorf(": %w", err)
				}
			}
		}
	}
	return
}
