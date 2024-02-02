// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_user

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func UserResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"access_rule": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"allow": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Description:         "List of access rule entries in the form `<verb>:<resource specifier>[:<SLA>]` that specify requests to allow",
						MarkdownDescription: "List of access rule entries in the form `<verb>:<resource specifier>[:<SLA>]` that specify requests to allow",
					},
					"deny": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Description:         "List of access rule entries in the form `<verb>:<resource specifier>` that specify requests to deny",
						MarkdownDescription: "List of access rule entries in the form `<verb>:<resource specifier>` that specify requests to deny",
					},
				},
				CustomType: AccessRuleType{
					ObjectType: types.ObjectType{
						AttrTypes: AccessRuleValue{}.AttributeTypes(ctx),
					},
				},
				Required:            true,
				Description:         "The rule specifying access for the user",
				MarkdownDescription: "The rule specifying access for the user",
			},
			"labels": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "User-defined labels attached to the resource that can be used for filtering",
				MarkdownDescription: "User-defined labels attached to the resource that can be used for filtering",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("[a-z][a-z0-9]*"), ""),
				},
			},
			"organization": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("[a-z][a-z0-9]*"), ""),
				},
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The password for the user",
				MarkdownDescription: "The password for the user",
			},
			"resource_version": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The version of the resource. When specified in a `PUT` request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.",
				MarkdownDescription: "The version of the resource. When specified in a `PUT` request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.",
			},
		},
	}
}

type UserModel struct {
	AccessRule      AccessRuleValue `tfsdk:"access_rule"`
	Labels          types.Map       `tfsdk:"labels"`
	Name            types.String    `tfsdk:"name"`
	Organization    types.String    `tfsdk:"organization"`
	Password        types.String    `tfsdk:"password"`
	ResourceVersion types.String    `tfsdk:"resource_version"`
}

var _ basetypes.ObjectTypable = AccessRuleType{}

type AccessRuleType struct {
	basetypes.ObjectType
}

func (t AccessRuleType) Equal(o attr.Type) bool {
	other, ok := o.(AccessRuleType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t AccessRuleType) String() string {
	return "AccessRuleType"
}

func (t AccessRuleType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	allowAttribute, ok := attributes["allow"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allow is missing from object`)

		return nil, diags
	}

	allowVal, ok := allowAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allow expected to be basetypes.ListValue, was: %T`, allowAttribute))
	}

	denyAttribute, ok := attributes["deny"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`deny is missing from object`)

		return nil, diags
	}

	denyVal, ok := denyAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`deny expected to be basetypes.ListValue, was: %T`, denyAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return AccessRuleValue{
		Allow: allowVal,
		Deny:  denyVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewAccessRuleValueNull() AccessRuleValue {
	return AccessRuleValue{
		state: attr.ValueStateNull,
	}
}

func NewAccessRuleValueUnknown() AccessRuleValue {
	return AccessRuleValue{
		state: attr.ValueStateUnknown,
	}
}

func NewAccessRuleValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (AccessRuleValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing AccessRuleValue Attribute Value",
				"While creating a AccessRuleValue value, a missing attribute value was detected. "+
					"A AccessRuleValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AccessRuleValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid AccessRuleValue Attribute Type",
				"While creating a AccessRuleValue value, an invalid attribute value was detected. "+
					"A AccessRuleValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AccessRuleValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("AccessRuleValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra AccessRuleValue Attribute Value",
				"While creating a AccessRuleValue value, an extra attribute value was detected. "+
					"A AccessRuleValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra AccessRuleValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewAccessRuleValueUnknown(), diags
	}

	allowAttribute, ok := attributes["allow"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allow is missing from object`)

		return NewAccessRuleValueUnknown(), diags
	}

	allowVal, ok := allowAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allow expected to be basetypes.ListValue, was: %T`, allowAttribute))
	}

	denyAttribute, ok := attributes["deny"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`deny is missing from object`)

		return NewAccessRuleValueUnknown(), diags
	}

	denyVal, ok := denyAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`deny expected to be basetypes.ListValue, was: %T`, denyAttribute))
	}

	if diags.HasError() {
		return NewAccessRuleValueUnknown(), diags
	}

	return AccessRuleValue{
		Allow: allowVal,
		Deny:  denyVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewAccessRuleValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) AccessRuleValue {
	object, diags := NewAccessRuleValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewAccessRuleValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t AccessRuleType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewAccessRuleValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewAccessRuleValueUnknown(), nil
	}

	if in.IsNull() {
		return NewAccessRuleValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewAccessRuleValueMust(AccessRuleValue{}.AttributeTypes(ctx), attributes), nil
}

func (t AccessRuleType) ValueType(ctx context.Context) attr.Value {
	return AccessRuleValue{}
}

var _ basetypes.ObjectValuable = AccessRuleValue{}

type AccessRuleValue struct {
	Allow basetypes.ListValue `tfsdk:"allow"`
	Deny  basetypes.ListValue `tfsdk:"deny"`
	state attr.ValueState
}

func (v AccessRuleValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["allow"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["deny"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Allow.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["allow"] = val

		val, err = v.Deny.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["deny"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v AccessRuleValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v AccessRuleValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v AccessRuleValue) String() string {
	return "AccessRuleValue"
}

func (v AccessRuleValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	allowVal, d := types.ListValue(types.StringType, v.Allow.Elements())

	diags.Append(d...)

	if d.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			"allow": basetypes.ListType{
				ElemType: types.StringType,
			},
			"deny": basetypes.ListType{
				ElemType: types.StringType,
			},
		}), diags
	}

	denyVal, d := types.ListValue(types.StringType, v.Deny.Elements())

	diags.Append(d...)

	if d.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			"allow": basetypes.ListType{
				ElemType: types.StringType,
			},
			"deny": basetypes.ListType{
				ElemType: types.StringType,
			},
		}), diags
	}

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"allow": basetypes.ListType{
				ElemType: types.StringType,
			},
			"deny": basetypes.ListType{
				ElemType: types.StringType,
			},
		},
		map[string]attr.Value{
			"allow": allowVal,
			"deny":  denyVal,
		})

	return objVal, diags
}

func (v AccessRuleValue) Equal(o attr.Value) bool {
	other, ok := o.(AccessRuleValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Allow.Equal(other.Allow) {
		return false
	}

	if !v.Deny.Equal(other.Deny) {
		return false
	}

	return true
}

func (v AccessRuleValue) Type(ctx context.Context) attr.Type {
	return AccessRuleType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v AccessRuleValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"allow": basetypes.ListType{
			ElemType: types.StringType,
		},
		"deny": basetypes.ListType{
			ElemType: types.StringType,
		},
	}
}