// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package framework

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasource "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
)

func getSchemas() (openapi3.Schemas, error) {
	swagger, err := openapi.GetSwagger()
	if err != nil {
		return nil, err
	}

	if swagger.Components != nil {
		return swagger.Components.Schemas, nil
	}
	return nil, nil
}

type SchemaOverride = func(*openapi3.Schema)

// WithDescription returns a SchemaOverride that overrides the description of
// the property at the specified path.
func WithDescription(path, description string) SchemaOverride {
	return func(oas *openapi3.Schema) {
		// Split the path into name and remaining
		parts := strings.SplitN(path, ".", 2)
		property, ok := oas.Properties[parts[0]]
		if ok && property != nil && property.Value != nil {
			if len(parts) == 1 {
				// Property has been found, so override its description
				property.Value.Description = description
			} else {
				// Property is nested, so invoke override on nested schema
				WithDescription(parts[1], description)(property.Value)
			}
		}
	}
}

func GetSchema(name string, overrides ...SchemaOverride) (*openapi3.Schema, error) {
	schemas, err := getSchemas()
	if err != nil {
		return nil, err
	}
	schemaRef, ok := schemas[name]
	if ok && schemaRef != nil && schemaRef.Value != nil {
		// Invoke schema overrides
		for _, override := range overrides {
			override(schemaRef.Value)
		}
		return schemaRef.Value, nil
	}
	return nil, fmt.Errorf("Schema %s not found", name)
}

func GetResourceAttributes(name string, overrides ...SchemaOverride) (map[string]resource.Attribute, error) {
	oas, err := GetSchema(name, overrides...)
	if err != nil {
		return nil, err
	}
	return ToResourceSchema(oas, false), nil
}

func GetDataSourceAttributes(name string, overrides ...SchemaOverride) (map[string]datasource.Attribute, error) {
	oas, err := GetSchema(name, overrides...)
	if err != nil {
		return nil, err
	}
	return ToDataSourceSchema(oas), nil
}

func GetAttributeName(oas *openapi3.Schema) string {
	if value, ok := oas.Extensions["x-tf-name"]; ok {
		strvalue, ok := value.(string)
		if ok {
			return strvalue
		}
	}
	return ""
}

func IsExtensionSet(oas *openapi3.Schema, key string) bool {
	if value, ok := oas.Extensions[key]; ok {
		switch v := value.(type) {
		case bool:
			return v
		case string:
			return v == "true"
		}
	}
	return false
}

func IsIdentifierAttribute(oas *openapi3.Schema) bool {
	return IsExtensionSet(oas, "x-tf-identifier")
}

func IsSensitiveAttribute(oas *openapi3.Schema) bool {
	return IsExtensionSet(oas, "x-tf-sensitive")
}

func IsImmutableAttribute(oas *openapi3.Schema) bool {
	return IsExtensionSet(oas, "x-immutable") || IsIdentifierAttribute(oas)
}

func GetStringValidators(oas *openapi3.Schema) []validator.String {
	var validators []validator.String

	format := oas.Format
	if format == "date-time" {
		timestampPattern := "^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$"
		validators = append(validators, stringvalidator.RegexMatches(regexp.MustCompile(timestampPattern), "must be a RFC 3339 date-time without fractional seconds"))
	}

	pattern := oas.Pattern
	if pattern != "" {
		if !strings.HasPrefix(pattern, "^") {
			pattern = "^" + pattern
		}
		if !strings.HasSuffix(pattern, "$") {
			pattern = pattern + "$"
		}
		validators = append(validators, stringvalidator.RegexMatches(regexp.MustCompile(pattern), "must match pattern: "+pattern))
	}
	return validators
}

// GetTerraformType returns the primitive type appearing in the supplied schema.
// Types "array" and "object" are ignored and should be handled by using
// ToResourceSchema() or ToDataSourceSchema().
func GetTerraformType(schemaRef *openapi3.SchemaRef) attr.Type {
	if schemaRef != nil && schemaRef.Value != nil {
		switch schemaRef.Value.Type {
		case "boolean":
			return types.BoolType
		case "integer":
			return types.Int64Type
		case "number":
			return types.NumberType
		case "string":
			return types.StringType
		}
	}
	return nil
}

func appendNonNil[T any](arr []T, elems ...T) []T {
	for _, elem := range elems {
		if !reflect.ValueOf(elem).IsNil() {
			arr = append(arr, elem)
		}
	}
	return arr
}

func ToResourceSchema(oas *openapi3.Schema, readOnly bool) map[string]resource.Attribute {
	if oas == nil {
		return nil
	}
	// Create set of required attributes
	required := make(map[string]struct{})
	for _, name := range oas.Required {
		required[name] = struct{}{}
	}
	// Convert JSONSchema properties to Terraform attributes
	attributes := make(map[string]resource.Attribute)
	for name, schema := range oas.Properties {
		// Skip schema if nil
		if schema == nil {
			continue
		}
		// Supply required value, which is attached to parent schema
		_, ok := required[name]
		tfname, tfschema := ToResourceAttribute(schema.Value, ok, readOnly)
		// If non-nil, then attribute should be exposed as resource attribute
		if tfschema != nil {
			attributes[tfname] = tfschema
		}
	}
	return attributes
}

func ToResourceAttribute(oas *openapi3.Schema, required, readOnly bool) (string, resource.Attribute) {
	if oas == nil {
		return "", nil
	}
	name := GetAttributeName(oas)
	if name == "" {
		return "", nil
	}
	if oas.ReadOnly {
		readOnly = oas.ReadOnly
	}
	var optional, computed bool
	if readOnly {
		required = false
		optional = false
		computed = true
	} else if !required {
		optional = true
		computed = true
	}
	sensitive := IsSensitiveAttribute(oas)
	// Add UseStateForUnknown plan modifier for computed attributes
	var useStateForUnknown *GenericPlanModifier
	if optional {
		useStateForUnknown = UseStateForUnknown()
	}
	// Add RequiresReplace plan modifier for immutable attributes
	var requiresReplace *GenericPlanModifier
	if IsImmutableAttribute(oas) {
		requiresReplace = RequiresReplace()
	}
	// Here comes the code duplication, required in order to interact with
	// the Terraform API...
	switch oas.Type {
	case "array":
		// If array contains objects, use ListNestedAttribute to attach nested object schema
		if oas.Items.Value != nil && oas.Items.Value.Type == "object" {
			return name, &resource.ListNestedAttribute{
				Description:         oas.Description,
				MarkdownDescription: oas.Description,
				Required:            required,
				Optional:            optional,
				Computed:            computed,
				Sensitive:           sensitive,
				PlanModifiers:       appendNonNil([]planmodifier.List{}, planmodifier.List(useStateForUnknown), planmodifier.List(requiresReplace)),
				NestedObject: resource.NestedAttributeObject{
					Attributes: ToResourceSchema(oas.Items.Value, readOnly),
				},
			}
		}
		return name, &resource.ListAttribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           sensitive,
			PlanModifiers:       appendNonNil([]planmodifier.List{}, planmodifier.List(useStateForUnknown), planmodifier.List(requiresReplace)),
			ElementType:         GetTerraformType(oas.Items),
		}
	case "boolean":
		return name, &resource.BoolAttribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           sensitive,
			PlanModifiers:       appendNonNil([]planmodifier.Bool{}, planmodifier.Bool(useStateForUnknown), planmodifier.Bool(requiresReplace)),
		}
	case "integer":
		return name, &resource.Int64Attribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           sensitive,
			PlanModifiers:       appendNonNil([]planmodifier.Int64{}, planmodifier.Int64(useStateForUnknown), planmodifier.Int64(requiresReplace)),
		}
	case "object":
		if oas.AdditionalProperties.Schema != nil {
			// If map values are objects, use MapNestedAttribute to attach nested object schema
			if oas.AdditionalProperties.Schema.Value != nil && oas.AdditionalProperties.Schema.Value.Type == "object" {
				return name, &resource.MapNestedAttribute{
					Description:         oas.Description,
					MarkdownDescription: oas.Description,
					Required:            required,
					Optional:            optional,
					Computed:            computed,
					Sensitive:           sensitive,
					NestedObject: resource.NestedAttributeObject{
						Attributes: ToResourceSchema(oas.AdditionalProperties.Schema.Value, readOnly),
					},
					PlanModifiers: appendNonNil([]planmodifier.Map{}, planmodifier.Map(useStateForUnknown), planmodifier.Map(requiresReplace)),
				}
			}
			return name, &resource.MapAttribute{
				Description:         oas.Description,
				MarkdownDescription: oas.Description,
				Required:            required,
				Optional:            optional,
				Computed:            computed,
				Sensitive:           sensitive,
				ElementType:         GetTerraformType(oas.AdditionalProperties.Schema),
				PlanModifiers:       appendNonNil([]planmodifier.Map{}, planmodifier.Map(useStateForUnknown), planmodifier.Map(requiresReplace)),
			}
		} else {
			return name, &resource.SingleNestedAttribute{
				Description:         oas.Description,
				MarkdownDescription: oas.Description,
				Required:            required,
				Optional:            optional,
				Computed:            computed,
				Sensitive:           sensitive,
				Attributes:          ToResourceSchema(oas, readOnly),
				PlanModifiers:       appendNonNil([]planmodifier.Object{}, planmodifier.Object(useStateForUnknown), planmodifier.Object(requiresReplace)),
			}
		}
	case "string":
		return name, &resource.StringAttribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           sensitive,
			Validators:          GetStringValidators(oas),
			PlanModifiers:       appendNonNil([]planmodifier.String{}, planmodifier.String(useStateForUnknown), planmodifier.String(requiresReplace)),
		}
	default:
		return "", nil
	}
}

func ToDataSourceSchema(oas *openapi3.Schema) map[string]datasource.Attribute {
	if oas == nil {
		return nil
	}
	// Convert JSONSchema properties to Terraform attributes
	attributes := make(map[string]datasource.Attribute)
	for _, schema := range oas.Properties {
		// Skip schema if nil
		if schema == nil {
			continue
		}
		tfname, tfschema := ToDataSourceAttribute(schema.Value)
		// If non-nil, then attribute should be exposed as datasource attribute
		if tfschema != nil {
			attributes[tfname] = tfschema
		}
	}
	return attributes
}

func ToDataSourceAttribute(oas *openapi3.Schema) (string, datasource.Attribute) {
	if oas == nil {
		return "", nil
	}
	name := GetAttributeName(oas)
	if name == "" {
		return "", nil
	}
	required := IsIdentifierAttribute(oas)
	computed := !required
	sensitive := IsSensitiveAttribute(oas)
	// Here comes the code duplication, required in order to interact with
	// the Terraform API...
	switch oas.Type {
	case "array":
		// If array contains objects, use ListNestedAttribute to attach nested object schema
		if oas.Items.Value != nil && oas.Items.Value.Type == "object" {
			return name, &datasource.ListNestedAttribute{
				Description:         oas.Description,
				MarkdownDescription: oas.Description,
				Required:            required,
				Computed:            computed,
				Sensitive:           sensitive,
				NestedObject: datasource.NestedAttributeObject{
					Attributes: ToDataSourceSchema(oas.Items.Value),
				},
			}
		}
		return name, &datasource.ListAttribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Computed:            computed,
			Sensitive:           sensitive,
			ElementType:         GetTerraformType(oas.Items),
		}
	case "boolean":
		return name, &datasource.BoolAttribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Computed:            computed,
			Sensitive:           sensitive,
		}
	case "integer":
		return name, &datasource.Int64Attribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Computed:            computed,
			Sensitive:           sensitive,
		}
	case "object":
		if oas.AdditionalProperties.Schema != nil {
			// If map values are objects, use MapNestedAttribute to attach nested object schema
			if oas.AdditionalProperties.Schema.Value != nil && oas.AdditionalProperties.Schema.Value.Type == "object" {
				return name, &datasource.MapNestedAttribute{
					Description:         oas.Description,
					MarkdownDescription: oas.Description,
					Required:            required,
					Computed:            computed,
					Sensitive:           sensitive,
					NestedObject: datasource.NestedAttributeObject{
						Attributes: ToDataSourceSchema(oas.AdditionalProperties.Schema.Value),
					},
				}
			}
			return name, &datasource.MapAttribute{
				Description:         oas.Description,
				MarkdownDescription: oas.Description,
				Required:            required,
				Computed:            computed,
				Sensitive:           sensitive,
				ElementType:         GetTerraformType(oas.AdditionalProperties.Schema),
			}
		} else {
			return name, &datasource.SingleNestedAttribute{
				Description:         oas.Description,
				MarkdownDescription: oas.Description,
				Required:            required,
				Computed:            computed,
				Sensitive:           sensitive,
				Attributes:          ToDataSourceSchema(oas),
			}
		}
	case "string":
		return name, &datasource.StringAttribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Computed:            computed,
			Sensitive:           sensitive,
		}
	default:
		return "", nil
	}
}
