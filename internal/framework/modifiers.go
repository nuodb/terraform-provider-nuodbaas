package framework

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ planmodifier.Bool    = GenericPlanModifier{}
	_ planmodifier.Int64   = GenericPlanModifier{}
	_ planmodifier.Float64 = GenericPlanModifier{}
	_ planmodifier.Number  = GenericPlanModifier{}
	_ planmodifier.String  = GenericPlanModifier{}
	_ planmodifier.List    = GenericPlanModifier{}
	_ planmodifier.Map     = GenericPlanModifier{}
	_ planmodifier.Object  = GenericPlanModifier{}
)

// GenericRequest is a request for plan modification that uses attr.Value
// instead of specific type of the attribute to enable code reuse across types.
type GenericRequest struct {
	// Path contains the path of the attribute for modification. Use this path
	// for any response diagnostics.
	Path path.Path

	// Config contains the entire configuration of the resource.
	Config tfsdk.Config

	// ConfigValue contains the value of the attribute for modification from the configuration.
	ConfigValue attr.Value

	// Plan contains the entire proposed new state of the resource.
	Plan tfsdk.Plan

	// PlanValue contains the value of the attribute for modification from the proposed new state.
	PlanValue attr.Value

	// State contains the entire prior state of the resource.
	State tfsdk.State

	// StateValue contains the value of the attribute for modification from the prior state.
	StateValue attr.Value
}

// GenericResponse is a response to a GenericRequest.
type GenericResponse struct {
	// PlanValue is the planned new state for the attribute.
	PlanValue attr.Value

	// RequiresReplace indicates whether a change in the attribute
	// requires replacement of the whole resource.
	RequiresReplace bool

	// Diagnostics report errors or warnings related to validating the data
	// source configuration. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics *diag.Diagnostics
}

// GenericPlanModifier is a plan modifier implementation that can be applied to
// any type.
type GenericPlanModifier struct {
	description string
	fn          func(req GenericRequest, resp *GenericResponse)
}

// Description returns a human-readable description of the plan modifier.
func (m GenericPlanModifier) Description(_ context.Context) string {
	return m.description
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m GenericPlanModifier) MarkdownDescription(_ context.Context) string {
	return m.description
}

func (m GenericPlanModifier) createRequest(
	attrPath path.Path,
	config tfsdk.Config,
	configValue attr.Value,
	plan tfsdk.Plan,
	planValue attr.Value,
	state tfsdk.State,
	stateValue attr.Value,
) GenericRequest {
	return GenericRequest{
		Path:        attrPath,
		Config:      config,
		ConfigValue: configValue,
		Plan:        plan,
		PlanValue:   planValue,
		State:       state,
		StateValue:  stateValue,
	}
}

func (m GenericPlanModifier) createResponse(
	planValue attr.Value,
	requiresReplace bool,
	diagnostics *diag.Diagnostics,
) GenericResponse {
	return GenericResponse{
		PlanValue:       planValue,
		RequiresReplace: requiresReplace,
		Diagnostics:     diagnostics,
	}
}

func (m GenericPlanModifier) execute(req GenericRequest, resp *GenericResponse) {
	if m.fn != nil {
		m.fn(req, resp)
	}
}

// PlanModifyBool implements the plan modification logic.
func (m GenericPlanModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	genericRequest := m.createRequest(req.Path, req.Config, req.ConfigValue, req.Plan, req.PlanValue, req.State, req.StateValue)
	genericResponse := m.createResponse(resp.PlanValue, resp.RequiresReplace, &resp.Diagnostics)
	m.execute(genericRequest, &genericResponse)
	resp.RequiresReplace = genericResponse.RequiresReplace
	resp.PlanValue = genericResponse.PlanValue.(types.Bool)
}

// PlanModifyInt64 implements the plan modification logic.
func (m GenericPlanModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	genericRequest := m.createRequest(req.Path, req.Config, req.ConfigValue, req.Plan, req.PlanValue, req.State, req.StateValue)
	genericResponse := m.createResponse(resp.PlanValue, resp.RequiresReplace, &resp.Diagnostics)
	m.execute(genericRequest, &genericResponse)
	resp.RequiresReplace = genericResponse.RequiresReplace
	resp.PlanValue = genericResponse.PlanValue.(types.Int64)
}

// PlanModifyFloat64 implements the plan modification logic.
func (m GenericPlanModifier) PlanModifyFloat64(_ context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	genericRequest := m.createRequest(req.Path, req.Config, req.ConfigValue, req.Plan, req.PlanValue, req.State, req.StateValue)
	genericResponse := m.createResponse(resp.PlanValue, resp.RequiresReplace, &resp.Diagnostics)
	m.execute(genericRequest, &genericResponse)
	resp.RequiresReplace = genericResponse.RequiresReplace
	resp.PlanValue = genericResponse.PlanValue.(types.Float64)
}

// PlanModifyNumber implements the plan modification logic.
func (m GenericPlanModifier) PlanModifyNumber(_ context.Context, req planmodifier.NumberRequest, resp *planmodifier.NumberResponse) {
	genericRequest := m.createRequest(req.Path, req.Config, req.ConfigValue, req.Plan, req.PlanValue, req.State, req.StateValue)
	genericResponse := m.createResponse(resp.PlanValue, resp.RequiresReplace, &resp.Diagnostics)
	m.execute(genericRequest, &genericResponse)
	resp.RequiresReplace = genericResponse.RequiresReplace
	resp.PlanValue = genericResponse.PlanValue.(types.Number)
}

// PlanModifyString implements the plan modification logic.
func (m GenericPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	genericRequest := m.createRequest(req.Path, req.Config, req.ConfigValue, req.Plan, req.PlanValue, req.State, req.StateValue)
	genericResponse := m.createResponse(resp.PlanValue, resp.RequiresReplace, &resp.Diagnostics)
	m.execute(genericRequest, &genericResponse)
	resp.RequiresReplace = genericResponse.RequiresReplace
	resp.PlanValue = genericResponse.PlanValue.(types.String)
}

// PlanModifyList implements the plan modification logic.
func (m GenericPlanModifier) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	genericRequest := m.createRequest(req.Path, req.Config, req.ConfigValue, req.Plan, req.PlanValue, req.State, req.StateValue)
	genericResponse := m.createResponse(resp.PlanValue, resp.RequiresReplace, &resp.Diagnostics)
	m.execute(genericRequest, &genericResponse)
	resp.RequiresReplace = genericResponse.RequiresReplace
	resp.PlanValue = genericResponse.PlanValue.(types.List)
}

// PlanModifyMap implements the plan modification logic.
func (m GenericPlanModifier) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	genericRequest := m.createRequest(req.Path, req.Config, req.ConfigValue, req.Plan, req.PlanValue, req.State, req.StateValue)
	genericResponse := m.createResponse(resp.PlanValue, resp.RequiresReplace, &resp.Diagnostics)
	m.execute(genericRequest, &genericResponse)
	resp.RequiresReplace = genericResponse.RequiresReplace
	resp.PlanValue = genericResponse.PlanValue.(types.Map)
}

// PlanModifyObject implements the plan modification logic.
func (m GenericPlanModifier) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	genericRequest := m.createRequest(req.Path, req.Config, req.ConfigValue, req.Plan, req.PlanValue, req.State, req.StateValue)
	genericResponse := m.createResponse(resp.PlanValue, resp.RequiresReplace, &resp.Diagnostics)
	m.execute(genericRequest, &genericResponse)
	resp.RequiresReplace = genericResponse.RequiresReplace
	resp.PlanValue = genericResponse.PlanValue.(types.Object)
}

func useStateForUnknown(req GenericRequest, resp *GenericResponse) {
	// Do nothing if there is no state. NOTE: This condition is different
	// from the one used by the standard UseStateForUnknown modifiers in the
	// Terraform library, which check whether the value of the specific
	// attribute is null. That check triggers unnecessary updates where the
	// optional/computed attribute had no default value generated for it by
	// the system, because the fact that an empty/null value was generated
	// by the system is not detected.
	if req.State.Raw.IsNull() {
		return
	}
	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}
	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}
	resp.PlanValue = req.StateValue
}

const (
	ALLOW_DESTRUCTIVE_REPLACE_VAR      = "NUODB_CP_ALLOW_DESTRUCTIVE_REPLACE"
	IMMUTABLE_ATTRIBUTE_CHANGE_ERR_FMT = "A change has been made to the immutable attribute `%s`, which may be rejected by the server or ignored. " +
		"In order for an immutable attribute change to take effect, it is necessary to delete and re-create the resource, which may result in data loss.\n\n" +
		"If it is acceptable to destroy all data associated with the resource, then deletion and re-creation can be enabled by setting the environment variable `" + ALLOW_DESTRUCTIVE_REPLACE_VAR + "=true` when running `terraform apply`:\n\n" +
		"    " + ALLOW_DESTRUCTIVE_REPLACE_VAR + "=true terraform apply\n\n" +
		"If the resource has dependent resources, it may be necessary to explicitly destroy and re-create the resource and its dependent resources by running the following:\n\n" +
		"    terraform destroy -target=<resource>\n    terraform apply\n\n" +
		"If it is not acceptable to destroy all data associated with the resource, the configured value of the attribute can be reverted to the value in the Terraform state (`%s`)."
)

func requiresReplace(req GenericRequest, resp *GenericResponse) {
	// Do not replace on resource creation.
	if req.State.Raw.IsNull() {
		return
	}
	// Do not replace on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}
	// Do not replace if the plan and state values are equal.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}
	if os.Getenv(ALLOW_DESTRUCTIVE_REPLACE_VAR) == "true" {
		resp.RequiresReplace = true
	} else {
		// Display a warning about a change to an immutable attribute,
		// which should fail to be enforced by the server. If the
		// attribute is not sensitive, then include the attribute value
		// in the output so that the user can set the configured value
		// to it to suppress the errors generated by the server.
		attributeValue := "<redacted>"
		attribute, err := req.State.Schema.AttributeAtPath(context.Background(), req.Path)
		if err == nil && !attribute.IsSensitive() {
			attributeValue = req.StateValue.String()
		}
		resp.Diagnostics.AddWarning("Immutable Attribute Change",
			fmt.Sprintf(IMMUTABLE_ATTRIBUTE_CHANGE_ERR_FMT, req.Path, attributeValue))
	}
}

func UseStateForUnknown() *GenericPlanModifier {
	return &GenericPlanModifier{
		description: "Once set, the value of this attribute in state will not change.",
		fn:          useStateForUnknown,
	}
}

func RequiresReplace() *GenericPlanModifier {
	return &GenericPlanModifier{
		description: "If the value of this attribute changes, Terraform will destroy and recreate the resource.",
		fn:          requiresReplace,
	}
}
