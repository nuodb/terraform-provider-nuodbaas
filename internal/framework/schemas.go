package framework

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasource "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
)

const (
	DatabaseSchemaName = "DatabaseCreateUpdateModel"
	ProjectSchemaName  = "ProjectModel"
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

func GetSchema(name string) (*openapi3.Schema, error) {
	schemas, err := getSchemas()
	if err != nil {
		return nil, err
	}
	schemaRef, ok := schemas[name]
	if ok || schemaRef != nil {
		if schemaRef.Value != nil {
			return schemaRef.Value, nil
		}
	}
	return nil, fmt.Errorf("Schema %s not found", name)
}

func GetDatabaseSchema() (*openapi3.Schema, error) {
	return GetSchema(DatabaseSchemaName)
}

func GetProjectSchema() (*openapi3.Schema, error) {
	return GetSchema(ProjectSchemaName)
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
	return IsExtensionSet(oas, "x-immutable")
}

func GetTerraformType(schemaRef *openapi3.SchemaRef) attr.Type {
	if schemaRef == nil {
		return nil
	}
	if schemaRef.Value == nil {
		return types.ObjectType{}
	}
	switch schemaRef.Value.Type {
	case "array":
		return types.ListType{}
	case "boolean":
		return types.BoolType
	case "integer":
		return types.Int64Type
	case "object":
		return types.ObjectType{}
	case "string":
		return types.StringType
	default:
		return nil
	}
}

func AppendNonNil[T any](arr []T, elems ...T) []T {
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
		return name, &resource.ListAttribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           sensitive,
			PlanModifiers:       AppendNonNil([]planmodifier.List{}, planmodifier.List(useStateForUnknown), planmodifier.List(requiresReplace)),
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
			PlanModifiers:       AppendNonNil([]planmodifier.Bool{}, planmodifier.Bool(useStateForUnknown), planmodifier.Bool(requiresReplace)),
		}
	case "integer":
		return name, &resource.Int64Attribute{
			Description:         oas.Description,
			MarkdownDescription: oas.Description,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           sensitive,
			PlanModifiers:       AppendNonNil([]planmodifier.Int64{}, planmodifier.Int64(useStateForUnknown), planmodifier.Int64(requiresReplace)),
		}
	case "object":
		if oas.AdditionalProperties.Schema != nil {
			return name, &resource.MapAttribute{
				Description:         oas.Description,
				MarkdownDescription: oas.Description,
				Required:            required,
				Optional:            optional,
				Computed:            computed,
				Sensitive:           sensitive,
				ElementType:         GetTerraformType(oas.AdditionalProperties.Schema),
				PlanModifiers:       AppendNonNil([]planmodifier.Map{}, planmodifier.Map(useStateForUnknown), planmodifier.Map(requiresReplace)),
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
				PlanModifiers:       AppendNonNil([]planmodifier.Object{}, planmodifier.Object(useStateForUnknown), planmodifier.Object(requiresReplace)),
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
			PlanModifiers:       AppendNonNil([]planmodifier.String{}, planmodifier.String(useStateForUnknown), planmodifier.String(requiresReplace)),
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
