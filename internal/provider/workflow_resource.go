package provider

import (
	"context"
	"fmt"

	"github.com/aws-gopher/terraform-provider-unstructured/internal/resource_workflow"
	"github.com/aws-gopher/unstructured-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*workflowResource)(nil)
var _ resource.ResourceWithConfigure = (*workflowResource)(nil)
var _ resource.ResourceWithImportState = (*workflowResource)(nil)

func NewWorkflowResource() resource.Resource {
	return &workflowResource{}
}

type workflowResource struct {
	client *unstructured.Client
}

func (r *workflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (r *workflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_workflow.WorkflowResourceSchema(ctx)
}

func (r *workflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*unstructured.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *unstructured.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *workflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_workflow.WorkflowModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert WorkflowNodes from Terraform model to API format
	var workflowNodes []unstructured.WorkflowNode
	if !data.WorkflowNodes.IsNull() && !data.WorkflowNodes.IsUnknown() {
		var nodeValues []resource_workflow.WorkflowNodesValue
		resp.Diagnostics.Append(data.WorkflowNodes.ElementsAs(ctx, &nodeValues, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		workflowNodes = make([]unstructured.WorkflowNode, len(nodeValues))
		for i, nodeValue := range nodeValues {
			// Convert settings from ObjectValue to map[string]interface{}
			var settings map[string]interface{}
			if !nodeValue.Settings.IsNull() && !nodeValue.Settings.IsUnknown() {
				settings = make(map[string]interface{})
				settingsMap := nodeValue.Settings.Attributes()
				for k, v := range settingsMap {
					if !v.IsNull() && !v.IsUnknown() {
						switch val := v.(type) {
						case types.String:
							settings[k] = val.ValueString()
						case types.Int64:
							settings[k] = val.ValueInt64()
						case types.Float64:
							settings[k] = val.ValueFloat64()
						case types.Bool:
							settings[k] = val.ValueBool()
						default:
							settings[k] = val.String()
						}
					}
				}
			}

			workflowNodes[i] = unstructured.WorkflowNode{
				ID:       nodeValue.Id.ValueStringPointer(),
				Name:     nodeValue.Name.ValueString(),
				Type:     nodeValue.WorkflowNodesType.ValueString(),
				Subtype:  nodeValue.Subtype.ValueString(),
				Settings: settings,
			}
		}
	}

	// Convert Schedule from string to API format
	var schedule *string
	if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
		scheduleValue := data.Schedule.ValueString()
		schedule = &scheduleValue
	}

	// Convert ReprocessAll
	var reprocessAll *bool
	if !data.ReprocessAll.IsNull() && !data.ReprocessAll.IsUnknown() {
		reprocessAllValue := data.ReprocessAll.ValueBool()
		reprocessAll = &reprocessAllValue
	}

	// Convert SourceID and DestinationID
	var sourceID *string
	var destinationID *string
	if !data.SourceId.IsNull() && !data.SourceId.IsUnknown() {
		sourceIDValue := data.SourceId.ValueString()
		sourceID = &sourceIDValue
	}
	if !data.DestinationId.IsNull() && !data.DestinationId.IsUnknown() {
		destinationIDValue := data.DestinationId.ValueString()
		destinationID = &destinationIDValue
	}

	// Create API call logic
	createRequest := unstructured.CreateWorkflowRequest{
		Name:          data.Name.ValueString(),
		SourceID:      sourceID,
		DestinationID: destinationID,
		WorkflowType:  unstructured.WorkflowType(data.WorkflowType.ValueString()),
		WorkflowNodes: workflowNodes,
		Schedule:      schedule,
		ReprocessAll:  reprocessAll,
	}

	workflow, err := r.client.CreateWorkflow(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating workflow", err.Error())
		return
	}

	// Convert the created workflow back to the model and set state
	resp.Diagnostics.Append(resp.State.Set(ctx, resource_workflow.WorkflowToModel(ctx, workflow, resp.Diagnostics))...)
}

func (r *workflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_workflow.WorkflowModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	workflow, err := r.client.GetWorkflow(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting workflow", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, resource_workflow.WorkflowToModel(ctx, workflow, resp.Diagnostics))...)
}

func (r *workflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_workflow.WorkflowModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state resource_workflow.WorkflowModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert WorkflowNodes from Terraform model to API format
	var workflowNodes []unstructured.WorkflowNode
	if !data.WorkflowNodes.IsNull() && !data.WorkflowNodes.IsUnknown() {
		var nodeValues []resource_workflow.WorkflowNodesValue
		resp.Diagnostics.Append(data.WorkflowNodes.ElementsAs(ctx, &nodeValues, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		workflowNodes = make([]unstructured.WorkflowNode, len(nodeValues))
		for i, nodeValue := range nodeValues {
			// Convert settings from ObjectValue to map[string]interface{}
			var settings map[string]interface{}
			if !nodeValue.Settings.IsNull() && !nodeValue.Settings.IsUnknown() {
				settings = make(map[string]interface{})
				settingsMap := nodeValue.Settings.Attributes()
				for k, v := range settingsMap {
					if !v.IsNull() && !v.IsUnknown() {
						switch val := v.(type) {
						case types.String:
							settings[k] = val.ValueString()
						case types.Int64:
							settings[k] = val.ValueInt64()
						case types.Float64:
							settings[k] = val.ValueFloat64()
						case types.Bool:
							settings[k] = val.ValueBool()
						default:
							settings[k] = val.String()
						}
					}
				}
			}

			workflowNodes[i] = unstructured.WorkflowNode{
				ID:       nodeValue.Id.ValueStringPointer(),
				Name:     nodeValue.Name.ValueString(),
				Type:     nodeValue.WorkflowNodesType.ValueString(),
				Subtype:  nodeValue.Subtype.ValueString(),
				Settings: settings,
			}
		}
	}

	// Convert Schedule from string to API format
	var schedule *string
	if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
		scheduleValue := data.Schedule.ValueString()
		schedule = &scheduleValue
	}

	// Convert ReprocessAll
	var reprocessAll *bool
	if !data.ReprocessAll.IsNull() && !data.ReprocessAll.IsUnknown() {
		reprocessAllValue := data.ReprocessAll.ValueBool()
		reprocessAll = &reprocessAllValue
	}

	// Convert SourceID and DestinationID
	var sourceID *string
	var destinationID *string
	if !data.SourceId.IsNull() && !data.SourceId.IsUnknown() {
		sourceIDValue := data.SourceId.ValueString()
		sourceID = &sourceIDValue
	}
	if !data.DestinationId.IsNull() && !data.DestinationId.IsUnknown() {
		destinationIDValue := data.DestinationId.ValueString()
		destinationID = &destinationIDValue
	}

	// Convert Name and WorkflowType to pointers for update
	var name *string
	var workflowType *unstructured.WorkflowType
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		nameValue := data.Name.ValueString()
		name = &nameValue
	}
	if !data.WorkflowType.IsNull() && !data.WorkflowType.IsUnknown() {
		workflowTypeValue := unstructured.WorkflowType(data.WorkflowType.ValueString())
		workflowType = &workflowTypeValue
	}

	// Create API call logic
	updateRequest := unstructured.UpdateWorkflowRequest{
		Name:          name,
		SourceID:      sourceID,
		DestinationID: destinationID,
		WorkflowType:  workflowType,
		WorkflowNodes: workflowNodes,
		Schedule:      schedule,
		ReprocessAll:  reprocessAll,
	}

	workflow, err := r.client.UpdateWorkflow(ctx, state.Id.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating workflow", err.Error())
		return
	}

	// Convert the updated workflow back to the model and set state
	resp.Diagnostics.Append(resp.State.Set(ctx, resource_workflow.WorkflowToModel(ctx, workflow, resp.Diagnostics))...)
}

func (r *workflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_workflow.WorkflowModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	err := r.client.DeleteWorkflow(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workflow", err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *workflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	workflow, err := r.client.GetWorkflow(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing workflow", err.Error())
		return
	}

	// Convert the workflow to the model and set state
	resp.Diagnostics.Append(resp.State.Set(ctx, resource_workflow.WorkflowToModel(ctx, workflow, resp.Diagnostics))...)
}
