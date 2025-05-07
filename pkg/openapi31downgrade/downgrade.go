package openapi31downgrade

import (
	"iter"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// DowngradeTo3_0 will attempt to downgrade the 3.1.x spec to v 3.0.3
//
// Only certain "transforms" are supported, and if the spec uses features that are not parsable/reresentable
// by [openapi3.T] then the code flow won't even reach this far.
func DowngradeTo3_0(spec *openapi3.T) (*openapi3.T, error) {
	if strings.HasPrefix(spec.OpenAPI, "3.0.") {
		return spec, nil
	}
	b, err := spec.MarshalJSON()
	if err != nil {
		return nil, err
	}

	doc, err := libopenapi.NewDocument(b)
	if err != nil {
		return nil, err
	}

	err = EditOpenAPISpec(doc)
	if err != nil {
		return nil, err
	}

	// render the document back to bytes and reload the model.
	blob, _, _, errors := doc.RenderAndReload()

	// if anything went wrong when re-rendering the v3 model, a slice of errors will be returned
	if len(errors) > 0 {
		err := (openapi3.MultiError)(errors)
		return nil, err
	}

	loader := openapi3.NewLoader()
	return loader.LoadFromData(blob)
}

// This _far_ from a complete conversion, it is just enough to get things working for the schema's I've
// used.
//
// More examples welcome!

func visitPaths(paths iter.Seq2[string, *v3.PathItem]) error {
	for path, item := range paths {
		for method, spec := range item.GetOperations().FromOldest() {
			loc := []any{path, method}
			for idx, param := range spec.Parameters {
				err := visitParam(param, append(loc, "parameters", idx))
				if err != nil {
					return err
				}
			}
			if spec.Responses.Default != nil {
				err := visitResponse(spec.Responses.Default, append(loc, "responses", "default"))
				if err != nil {
					return err
				}
			}
			for code, resp := range spec.Responses.Codes.FromNewest() {
				err := visitResponse(resp, append(loc, "responses", code))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func visitResponse(response *v3.Response, loc []any) error {
	for ct, media := range response.Content.FromOldest() {
		err := visitSchemaProxy(&media.Schema, append(loc, ct, "schema"))
		if err != nil {
			return err
		}
	}
	return nil
}

func visitComponents(components *v3.Components) error {
	for name, resp := range components.Responses.FromOldest() {
		err := visitResponse(resp, []any{"components", "responses", name})
		if err != nil {
			return err
		}
	}

	for name, proxy := range components.Schemas.FromOldest() {
		if proxy.IsReference() {
			continue
		}
		if schema := proxy.Schema(); schema != nil {
			newSchema, err := visitSchema(schema, []any{"components", "schemas", name})
			if err != nil {
				return err
			}
			if newSchema != nil {
				components.Schemas.Set(name, base.CreateSchemaProxy(newSchema))
			}
		}
	}

	return nil
}

func visitParam(param *v3.Parameter, loc []any) error {
	return visitSchemaProxy(&param.Schema, append(loc, "schema"))
}

func visitSchemaProxy(schemaProxy **base.SchemaProxy, loc []any) error {
	if (*schemaProxy).IsReference() {
		return nil
	}
	if schema := (*schemaProxy).Schema(); schema != nil {
		schema, err := visitSchema(schema, loc)
		if err != nil {
			return err
		}
		if schema != nil {
			(*schemaProxy) = base.CreateSchemaProxy(schema)
		}
	}
	return nil
}

func visitSchema(schema *base.Schema, loc []any) (*base.Schema, error) {
	changed := false
	if numOneOfs := len(schema.OneOf); numOneOfs > 0 {
		lastType := schema.OneOf[numOneOfs-1].Schema().Type
		if len(lastType) == 1 && lastType[0] == "null" {
			schema.OneOf = schema.OneOf[0 : numOneOfs-1]
			nullable := true
			numOneOfs--
			schema.Nullable = &nullable

			// If it was only a single type left, collapse it.
			if numOneOfs == 1 {
				oneOf := schema.OneOf[0].Schema()
				schema.OneOf = nil
				copyNonZeroProps(schema, oneOf)
			}
			changed = true
		}
	} else if len(schema.Type) == 1 && schema.Type[0] == "null" {
		schema.Type = nil
		nullable := true
		schema.Nullable = &nullable
		changed = true
	}

	for name, propProxy := range schema.Properties.FromOldest() {
		if propProxy.IsReference() {
			continue
		}
		lloc := append(loc, "properties", name)
		if prop := propProxy.Schema(); schema != nil {
			newSchema, err := visitSchema(prop, lloc)
			if err != nil {
				return nil, err
			}
			if newSchema != nil {
				schema.Properties.Set(name, base.CreateSchemaProxy(newSchema))
				changed = true
			}
		}
	}

	if changed {
		return schema, nil
	}
	return nil, nil
}

func copyNonZeroProps[T any](dest, src *T) {
	source := reflect.ValueOf(src)
	target := reflect.ValueOf(dest)
	for i := range source.Elem().NumField() {
		fld := source.Elem().Type().Field(i)
		if !fld.IsExported() {
			continue
		}
		val := source.Elem().Field(i)
		if !val.IsZero() {
			target.Elem().Field(i).Set(val)
		}
	}
}

func EditOpenAPISpec(doc libopenapi.Document) error {
	// because we know this is a v3 spec, we can build a ready to go model from it.
	v3Model, errors := doc.BuildV3Model()

	// if anything went wrong when building the v3 model, a slice of errors will be returned
	if len(errors) > 0 {
		return (openapi3.MultiError)(errors)
	}

	visitPaths(v3Model.Model.Paths.PathItems.FromOldest())
	visitComponents(v3Model.Model.Components)
	v3Model.Model.Version = "3.0.3"

	return nil
}
