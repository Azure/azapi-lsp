//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armalertsmanagement

import "time"

// Action to be applied.
type Action struct {
	// REQUIRED; Action that should be applied.
	ActionType *ActionType
}

// GetAction implements the ActionClassification interface for type Action.
func (a *Action) GetAction() *Action { return a }

// ActionGroup - A pointer to an Azure Action Group.
type ActionGroup struct {
	// REQUIRED; The resource ID of the Action Group. This cannot be null or empty.
	ActionGroupID *string

	// Predefined list of properties and configuration items for the action group.
	ActionProperties map[string]*string

	// the dictionary of custom properties to include with the post operation. These data are appended to the webhook payload.
	WebhookProperties map[string]*string
}

// ActionList - A list of Activity Log Alert rule actions.
type ActionList struct {
	// The list of the Action Groups.
	ActionGroups []*ActionGroup
}

// ActionStatus - Action status
type ActionStatus struct {
	// Value indicating whether alert is suppressed.
	IsSuppressed *bool
}

// AddActionGroups - Add action groups to alert processing rule.
type AddActionGroups struct {
	// REQUIRED; List of action group Ids to add to alert processing rule.
	ActionGroupIDs []*string

	// REQUIRED; Action that should be applied.
	ActionType *ActionType
}

// GetAction implements the ActionClassification interface for type AddActionGroups.
func (a *AddActionGroups) GetAction() *Action {
	return &Action{
		ActionType: a.ActionType,
	}
}

// Alert - An alert created in alert management service.
type Alert struct {
	// Alert property bag
	Properties *AlertProperties

	// READ-ONLY; Azure resource Id
	ID *string

	// READ-ONLY; Azure resource name
	Name *string

	// READ-ONLY; Azure resource type
	Type *string
}

// AlertModification - Alert Modification details
type AlertModification struct {
	// Properties of the alert modification item.
	Properties *AlertModificationProperties

	// READ-ONLY; Azure resource Id
	ID *string

	// READ-ONLY; Azure resource name
	Name *string

	// READ-ONLY; Azure resource type
	Type *string
}

// AlertModificationItem - Alert modification item.
type AlertModificationItem struct {
	// Modification comments
	Comments *string

	// Description of the modification
	Description *string

	// Reason for the modification
	ModificationEvent *AlertModificationEvent

	// Modified date and time
	ModifiedAt *string

	// Modified user details (Principal client name)
	ModifiedBy *string

	// New value
	NewValue *string

	// Old value
	OldValue *string
}

// AlertModificationProperties - Properties of the alert modification item.
type AlertModificationProperties struct {
	// Modification details
	Modifications []*AlertModificationItem

	// READ-ONLY; Unique Id of the alert for which the history is being retrieved
	AlertID *string
}

// AlertProcessingRule - Alert processing rule object containing target scopes, conditions and scheduling logic.
type AlertProcessingRule struct {
	// REQUIRED; Resource location
	Location *string

	// Alert processing rule properties.
	Properties *AlertProcessingRuleProperties

	// Resource tags
	Tags map[string]*string

	// READ-ONLY; Azure resource Id
	ID *string

	// READ-ONLY; Azure resource name
	Name *string

	// READ-ONLY; Alert processing rule system data.
	SystemData *SystemData

	// READ-ONLY; Azure resource type
	Type *string
}

// AlertProcessingRuleProperties - Alert processing rule properties defining scopes, conditions and scheduling logic for alert
// processing rule.
type AlertProcessingRuleProperties struct {
	// REQUIRED; Actions to be applied.
	Actions []ActionClassification

	// REQUIRED; Scopes on which alert processing rule will apply.
	Scopes []*string

	// Conditions on which alerts will be filtered.
	Conditions []*Condition

	// Description of alert processing rule.
	Description *string

	// Indicates if the given alert processing rule is enabled or disabled.
	Enabled *bool

	// Scheduling for alert processing rule.
	Schedule *Schedule
}

// AlertProcessingRulesList - List of alert processing rules.
type AlertProcessingRulesList struct {
	// URL to fetch the next set of alert processing rules.
	NextLink *string

	// List of alert processing rules.
	Value []*AlertProcessingRule
}

// AlertProperties - Alert property bag
type AlertProperties struct {
	// This object contains consistent fields across different monitor services.
	Essentials *Essentials

	// READ-ONLY; Information specific to the monitor service that gives more contextual details about the alert.
	Context any

	// READ-ONLY; Config which would be used for displaying the data in portal.
	EgressConfig any
}

// AlertRuleAllOfCondition - An Activity Log Alert rule condition that is met when all its member conditions are met.
type AlertRuleAllOfCondition struct {
	// REQUIRED; The list of Activity Log Alert rule conditions.
	AllOf []*AlertRuleAnyOfOrLeafCondition
}

// AlertRuleAnyOfOrLeafCondition - An Activity Log Alert rule condition that is met when all its member conditions are met.
// Each condition can be of one of the following types:Important: Each type has its unique subset of properties.
// Properties from different types CANNOT exist in one condition.
// * Leaf Condition - must contain 'field' and either 'equals' or 'containsAny'.Please note, 'anyOf' should not be set in
// a Leaf Condition.
// * AnyOf Condition - must contain only 'anyOf' (which is an array of Leaf Conditions).Please note, 'field', 'equals' and
// 'containsAny' should not be set in an AnyOf Condition.
type AlertRuleAnyOfOrLeafCondition struct {
	// An Activity Log Alert rule condition that is met when at least one of its member leaf conditions are met.
	AnyOf []*AlertRuleLeafCondition

	// The value of the event's field will be compared to the values in this array (case-insensitive) to determine if the condition
	// is met.
	ContainsAny []*string

	// The value of the event's field will be compared to this value (case-insensitive) to determine if the condition is met.
	Equals *string

	// The name of the Activity Log event's field that this condition will examine. The possible values for this field are (case-insensitive):
	// 'resourceId', 'category', 'caller', 'level', 'operationName',
	// 'resourceGroup', 'resourceProvider', 'status', 'subStatus', 'resourceType', or anything beginning with 'properties'.
	Field *string
}

// AlertRuleLeafCondition - An Activity Log Alert rule condition that is met by comparing the field and value of an Activity
// Log event. This condition must contain 'field' and either 'equals' or 'containsAny'.
type AlertRuleLeafCondition struct {
	// The value of the event's field will be compared to the values in this array (case-insensitive) to determine if the condition
	// is met.
	ContainsAny []*string

	// The value of the event's field will be compared to this value (case-insensitive) to determine if the condition is met.
	Equals *string

	// The name of the Activity Log event's field that this condition will examine. The possible values for this field are (case-insensitive):
	// 'resourceId', 'category', 'caller', 'level', 'operationName',
	// 'resourceGroup', 'resourceProvider', 'status', 'subStatus', 'resourceType', or anything beginning with 'properties'.
	Field *string
}

// AlertRuleProperties - An Azure Activity Log Alert rule.
type AlertRuleProperties struct {
	// REQUIRED; The actions that will activate when the condition is met.
	Actions *ActionList

	// REQUIRED; The condition that will cause this alert to activate.
	Condition *AlertRuleAllOfCondition

	// A description of this Activity Log Alert rule.
	Description *string

	// Indicates whether this Activity Log Alert rule is enabled. If an Activity Log Alert rule is not enabled, then none of its
	// actions will be activated.
	Enabled *bool

	// A list of resource IDs that will be used as prefixes. The alert will only apply to Activity Log events with resource IDs
	// that fall under one of these prefixes. This list must include at least one
	// item.
	Scopes []*string

	// The tenant GUID. Must be provided for tenant-level and management group events rules.
	TenantScope *string
}

// AlertRuleRecommendationProperties - Describes the format of Alert Rule Recommendations response.
type AlertRuleRecommendationProperties struct {
	// REQUIRED; The recommendation alert rule type.
	AlertRuleType *string

	// REQUIRED; A dictionary that provides the display information for an alert rule recommendation.
	DisplayInformation map[string]*string

	// REQUIRED; A complete ARM template to deploy the alert rules.
	RuleArmTemplate *RuleArmTemplate
}

// AlertRuleRecommendationResource - A single alert rule recommendation resource.
type AlertRuleRecommendationResource struct {
	// REQUIRED; recommendation properties.
	Properties *AlertRuleRecommendationProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// AlertRuleRecommendationsListResponse - List of alert rule recommendations.
type AlertRuleRecommendationsListResponse struct {
	// REQUIRED; the values for the alert rule recommendations.
	Value []*AlertRuleRecommendationResource

	// URL to fetch the next set of recommendations.
	NextLink *string
}

// AlertsList - List the alerts.
type AlertsList struct {
	// URL to fetch the next set of alerts.
	NextLink *string

	// List of alerts
	Value []*Alert
}

// AlertsMetaData - alert meta data information.
type AlertsMetaData struct {
	// alert meta data property bag
	Properties AlertsMetaDataPropertiesClassification
}

// AlertsMetaDataProperties - alert meta data property bag
type AlertsMetaDataProperties struct {
	// REQUIRED; Identification of the information to be retrieved by API call
	MetadataIdentifier *MetadataIdentifier
}

// GetAlertsMetaDataProperties implements the AlertsMetaDataPropertiesClassification interface for type AlertsMetaDataProperties.
func (a *AlertsMetaDataProperties) GetAlertsMetaDataProperties() *AlertsMetaDataProperties { return a }

// AlertsSummary - Summary of alerts based on the input filters and 'groupby' parameters.
type AlertsSummary struct {
	// Group the result set.
	Properties *AlertsSummaryGroup

	// READ-ONLY; Azure resource Id
	ID *string

	// READ-ONLY; Azure resource name
	Name *string

	// READ-ONLY; Azure resource type
	Type *string
}

// AlertsSummaryGroup - Group the result set.
type AlertsSummaryGroup struct {
	// Name of the field aggregated
	Groupedby *string

	// Total count of the smart groups.
	SmartGroupsCount *int64

	// Total count of the result set.
	Total *int64

	// List of the items
	Values []*AlertsSummaryGroupItem
}

// AlertsSummaryGroupItem - Alerts summary group item
type AlertsSummaryGroupItem struct {
	// Count of the aggregated field
	Count *int64

	// Name of the field aggregated
	Groupedby *string

	// Value of the aggregated field
	Name *string

	// List of the items
	Values []*AlertsSummaryGroupItem
}

// Comments - Change alert state reason
type Comments struct {
	Comments *string
}

// Condition to trigger an alert processing rule.
type Condition struct {
	// Field for a given condition.
	Field *Field

	// Operator for a given condition.
	Operator *Operator

	// List of values to match for a given condition.
	Values []*string
}

// DailyRecurrence - Daily recurrence object.
type DailyRecurrence struct {
	// REQUIRED; Specifies when the recurrence should be applied.
	RecurrenceType *RecurrenceType

	// End time for recurrence.
	EndTime *string

	// Start time for recurrence.
	StartTime *string
}

// GetRecurrence implements the RecurrenceClassification interface for type DailyRecurrence.
func (d *DailyRecurrence) GetRecurrence() *Recurrence {
	return &Recurrence{
		EndTime:        d.EndTime,
		RecurrenceType: d.RecurrenceType,
		StartTime:      d.StartTime,
	}
}

// Essentials - This object contains consistent fields across different monitor services.
type Essentials struct {
	// Action status
	ActionStatus *ActionStatus

	// Alert description.
	Description *string

	// Target ARM resource, on which alert got created.
	TargetResource *string

	// Resource group of target ARM resource, on which alert got created.
	TargetResourceGroup *string

	// Name of the target ARM resource name, on which alert got created.
	TargetResourceName *string

	// Resource type of target ARM resource, on which alert got created.
	TargetResourceType *string

	// READ-ONLY; Rule(monitor) which fired alert instance. Depending on the monitor service, this would be ARM id or name of
	// the rule.
	AlertRule *string

	// READ-ONLY; Alert object state, which can be modified by the user.
	AlertState *AlertState

	// READ-ONLY; Last modification time(ISO-8601 format) of alert instance.
	LastModifiedDateTime *time.Time

	// READ-ONLY; User who last modified the alert, in case of monitor service updates user would be 'system', otherwise name
	// of the user.
	LastModifiedUserName *string

	// READ-ONLY; Condition of the rule at the monitor service. It represents whether the underlying conditions have crossed the
	// defined alert rule thresholds.
	MonitorCondition *MonitorCondition

	// READ-ONLY; Resolved time(ISO-8601 format) of alert instance. This will be updated when monitor service resolves the alert
	// instance because the rule condition is no longer met.
	MonitorConditionResolvedDateTime *time.Time

	// READ-ONLY; Monitor service on which the rule(monitor) is set.
	MonitorService *MonitorService

	// READ-ONLY; Severity of alert Sev0 being highest and Sev4 being lowest.
	Severity *Severity

	// READ-ONLY; The type of signal the alert is based on, which could be metrics, logs or activity logs.
	SignalType *SignalType

	// READ-ONLY; Unique Id of the smart group
	SmartGroupID *string

	// READ-ONLY; Verbose reason describing the reason why this alert instance is added to a smart group
	SmartGroupingReason *string

	// READ-ONLY; Unique Id created by monitor service for each alert instance. This could be used to track the issue at the monitor
	// service, in case of Nagios, Zabbix, SCOM etc.
	SourceCreatedID *string

	// READ-ONLY; Creation time(ISO-8601 format) of alert instance.
	StartDateTime *time.Time
}

// MonitorServiceDetails - Details of a monitor service
type MonitorServiceDetails struct {
	// Monitor service display name
	DisplayName *string

	// Monitor service name
	Name *string
}

// MonitorServiceList - Monitor service details
type MonitorServiceList struct {
	// REQUIRED; Array of operations
	Data []*MonitorServiceDetails

	// REQUIRED; Identification of the information to be retrieved by API call
	MetadataIdentifier *MetadataIdentifier
}

// GetAlertsMetaDataProperties implements the AlertsMetaDataPropertiesClassification interface for type MonitorServiceList.
func (m *MonitorServiceList) GetAlertsMetaDataProperties() *AlertsMetaDataProperties {
	return &AlertsMetaDataProperties{
		MetadataIdentifier: m.MetadataIdentifier,
	}
}

// MonthlyRecurrence - Monthly recurrence object.
type MonthlyRecurrence struct {
	// REQUIRED; Specifies the values for monthly recurrence pattern.
	DaysOfMonth []*int32

	// REQUIRED; Specifies when the recurrence should be applied.
	RecurrenceType *RecurrenceType

	// End time for recurrence.
	EndTime *string

	// Start time for recurrence.
	StartTime *string
}

// GetRecurrence implements the RecurrenceClassification interface for type MonthlyRecurrence.
func (m *MonthlyRecurrence) GetRecurrence() *Recurrence {
	return &Recurrence{
		EndTime:        m.EndTime,
		RecurrenceType: m.RecurrenceType,
		StartTime:      m.StartTime,
	}
}

// Operation provided by provider
type Operation struct {
	// Properties of the operation
	Display *OperationDisplay

	// Name of the operation
	Name *string

	// Origin of the operation
	Origin *string
}

// OperationDisplay - Properties of the operation
type OperationDisplay struct {
	// Description of the operation
	Description *string

	// Operation name
	Operation *string

	// Provider name
	Provider *string

	// Resource name
	Resource *string
}

// OperationsList - Lists the operations available in the AlertsManagement RP.
type OperationsList struct {
	// REQUIRED; Array of operations
	Value []*Operation

	// URL to fetch the next set of alerts.
	NextLink *string
}

// PatchObject - Data contract for patch.
type PatchObject struct {
	// Properties supported by patch operation.
	Properties *PatchProperties

	// Tags to be updated.
	Tags map[string]*string
}

// PatchProperties - Alert processing rule properties supported by patch.
type PatchProperties struct {
	// Indicates if the given alert processing rule is enabled or disabled.
	Enabled *bool
}

type PrometheusRule struct {
	// REQUIRED; the expression to run for the rule.
	Expression *string

	// The array of actions that are performed when the alert rule becomes active, and when an alert condition is resolved. Only
	// relevant for alerts.
	Actions []*PrometheusRuleGroupAction

	// the name of the alert rule.
	Alert *string

	// annotations for rule group. Only relevant for alerts.
	Annotations map[string]*string

	// the flag that indicates whether the Prometheus rule is enabled.
	Enabled *bool

	// the amount of time alert must be active before firing. Only relevant for alerts.
	For *string

	// labels for rule group. Only relevant for alerts.
	Labels map[string]*string

	// the name of the recording rule.
	Record *string

	// defines the configuration for resolving fired alerts. Only relevant for alerts.
	ResolveConfiguration *PrometheusRuleResolveConfiguration

	// the severity of the alerts fired by the rule. Only relevant for alerts.
	Severity *int32
}

// PrometheusRuleGroupAction - An alert action. Only relevant for alerts.
type PrometheusRuleGroupAction struct {
	// The resource id of the action group to use.
	ActionGroupID *string

	// The properties of an action group object.
	ActionProperties map[string]*string
}

// PrometheusRuleGroupProperties - An alert rule.
type PrometheusRuleGroupProperties struct {
	// REQUIRED; defines the rules in the Prometheus rule group.
	Rules []*PrometheusRule

	// REQUIRED; the list of resource id's that this rule group is scoped to.
	Scopes []*string

	// the cluster name of the rule group evaluation.
	ClusterName *string

	// the description of the Prometheus rule group that will be included in the alert email.
	Description *string

	// the flag that indicates whether the Prometheus rule group is enabled.
	Enabled *bool

	// the interval in which to run the Prometheus rule group represented in ISO 8601 duration format. Should be between 1 and
	// 15 minutes
	Interval *string
}

// PrometheusRuleGroupResource - The Prometheus rule group resource.
type PrometheusRuleGroupResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// REQUIRED; The Prometheus rule group properties of the resource.
	Properties *PrometheusRuleGroupProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// PrometheusRuleGroupResourceCollection - Represents a collection of alert rule resources.
type PrometheusRuleGroupResourceCollection struct {
	// the values for the alert rule resources.
	Value []*PrometheusRuleGroupResource
}

// PrometheusRuleGroupResourcePatch - The Prometheus rule group resource for patch operations.
type PrometheusRuleGroupResourcePatch struct {
	Properties *PrometheusRuleGroupResourcePatchProperties

	// Resource tags
	Tags map[string]*string
}

type PrometheusRuleGroupResourcePatchProperties struct {
	// the flag that indicates whether the Prometheus rule group is enabled.
	Enabled *bool
}

// PrometheusRuleResolveConfiguration - Specifies the Prometheus alert rule configuration.
type PrometheusRuleResolveConfiguration struct {
	// the flag that indicates whether or not to auto resolve a fired alert.
	AutoResolved *bool

	// the duration a rule must evaluate as healthy before the fired alert is automatically resolved represented in ISO 8601 duration
	// format. Should be between 1 and 15 minutes
	TimeToResolve *string
}

// Recurrence object.
type Recurrence struct {
	// REQUIRED; Specifies when the recurrence should be applied.
	RecurrenceType *RecurrenceType

	// End time for recurrence.
	EndTime *string

	// Start time for recurrence.
	StartTime *string
}

// GetRecurrence implements the RecurrenceClassification interface for type Recurrence.
func (r *Recurrence) GetRecurrence() *Recurrence { return r }

// RemoveAllActionGroups - Indicates if all action groups should be removed.
type RemoveAllActionGroups struct {
	// REQUIRED; Action that should be applied.
	ActionType *ActionType
}

// GetAction implements the ActionClassification interface for type RemoveAllActionGroups.
func (r *RemoveAllActionGroups) GetAction() *Action {
	return &Action{
		ActionType: r.ActionType,
	}
}

// RuleArmTemplate - A complete ARM template to deploy the alert rules.
type RuleArmTemplate struct {
	// REQUIRED; A 4 number format for the version number of this template file. For example, 1.0.0.0
	ContentVersion *string

	// REQUIRED; Input parameter definitions
	Parameters any

	// REQUIRED; Alert rule resource definitions
	Resources []any

	// REQUIRED; JSON schema reference
	Schema *string

	// REQUIRED; Variable definitions
	Variables any
}

// Schedule - Scheduling configuration for a given alert processing rule.
type Schedule struct {
	// Scheduling effective from time. Date-Time in ISO-8601 format without timezone suffix.
	EffectiveFrom *string

	// Scheduling effective until time. Date-Time in ISO-8601 format without timezone suffix.
	EffectiveUntil *string

	// List of recurrences.
	Recurrences []RecurrenceClassification

	// Scheduling time zone.
	TimeZone *string
}

// SmartGroup - Set of related alerts grouped together smartly by AMS.
type SmartGroup struct {
	// Properties of smart group.
	Properties *SmartGroupProperties

	// READ-ONLY; Azure resource Id
	ID *string

	// READ-ONLY; Azure resource name
	Name *string

	// READ-ONLY; Azure resource type
	Type *string
}

// SmartGroupAggregatedProperty - Aggregated property of each type
type SmartGroupAggregatedProperty struct {
	// Total number of items of type.
	Count *int64

	// Name of the type.
	Name *string
}

// SmartGroupModification - Alert Modification details
type SmartGroupModification struct {
	// Properties of the smartGroup modification item.
	Properties *SmartGroupModificationProperties

	// READ-ONLY; Azure resource Id
	ID *string

	// READ-ONLY; Azure resource name
	Name *string

	// READ-ONLY; Azure resource type
	Type *string
}

// SmartGroupModificationItem - smartGroup modification item.
type SmartGroupModificationItem struct {
	// Modification comments
	Comments *string

	// Description of the modification
	Description *string

	// Reason for the modification
	ModificationEvent *SmartGroupModificationEvent

	// Modified date and time
	ModifiedAt *string

	// Modified user details (Principal client name)
	ModifiedBy *string

	// New value
	NewValue *string

	// Old value
	OldValue *string
}

// SmartGroupModificationProperties - Properties of the smartGroup modification item.
type SmartGroupModificationProperties struct {
	// Modification details
	Modifications []*SmartGroupModificationItem

	// URL to fetch the next set of results.
	NextLink *string

	// READ-ONLY; Unique Id of the smartGroup for which the history is being retrieved
	SmartGroupID *string
}

// SmartGroupProperties - Properties of smart group.
type SmartGroupProperties struct {
	// Summary of alertSeverities in the smart group
	AlertSeverities []*SmartGroupAggregatedProperty

	// Summary of alertStates in the smart group
	AlertStates []*SmartGroupAggregatedProperty

	// Total number of alerts in smart group
	AlertsCount *int64

	// Summary of monitorConditions in the smart group
	MonitorConditions []*SmartGroupAggregatedProperty

	// Summary of monitorServices in the smart group
	MonitorServices []*SmartGroupAggregatedProperty

	// The URI to fetch the next page of alerts. Call ListNext() with this URI to fetch the next page alerts.
	NextLink *string

	// Summary of target resource groups in the smart group
	ResourceGroups []*SmartGroupAggregatedProperty

	// Summary of target resource types in the smart group
	ResourceTypes []*SmartGroupAggregatedProperty

	// Summary of target resources in the smart group
	Resources []*SmartGroupAggregatedProperty

	// READ-ONLY; Last updated time of smart group. Date-Time in ISO-8601 format.
	LastModifiedDateTime *time.Time

	// READ-ONLY; Last modified by user name.
	LastModifiedUserName *string

	// READ-ONLY; Severity of smart group is the highest(Sev0 >… > Sev4) severity of all the alerts in the group.
	Severity *Severity

	// READ-ONLY; Smart group state
	SmartGroupState *State

	// READ-ONLY; Creation time of smart group. Date-Time in ISO-8601 format.
	StartDateTime *time.Time
}

// SmartGroupsList - List the alerts.
type SmartGroupsList struct {
	// URL to fetch the next set of alerts.
	NextLink *string

	// List of alerts
	Value []*SmartGroup
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time

	// The identity that created the resource.
	CreatedBy *string

	// The type of identity that created the resource.
	CreatedByType *CreatedByType

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time

	// The identity that last modified the resource.
	LastModifiedBy *string

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType
}

// TenantActivityLogAlertResource - A Tenant Activity Log Alert rule resource.
type TenantActivityLogAlertResource struct {
	// REQUIRED; The Activity Log Alert rule properties of the resource.
	Properties *AlertRuleProperties

	// The location of the resource. Since Azure Activity Log Alerts is a global service, the location of the rules should always
	// be 'global'.
	Location *string

	// The tags of the resource.
	Tags map[string]*string

	// READ-ONLY; The resource Id.
	ID *string

	// READ-ONLY; The name of the resource.
	Name *string

	// READ-ONLY; The type of the resource.
	Type *string
}

// TenantAlertRuleList - A list of Tenant Activity Log Alert rules.
type TenantAlertRuleList struct {
	// Provides the link to retrieve the next set of elements.
	NextLink *string

	// The list of Tenant Activity Log Alert rules.
	Value []*TenantActivityLogAlertResource
}

// TenantAlertRulePatchObject - An Activity Log Alert rule object for the body of patch operations.
type TenantAlertRulePatchObject struct {
	// The activity log alert settings for an update operation.
	Properties *TenantAlertRulePatchProperties

	// The resource tags
	Tags map[string]*string
}

// TenantAlertRulePatchProperties - An Activity Log Alert rule properties for patch operations.
type TenantAlertRulePatchProperties struct {
	// Indicates whether this Activity Log Alert rule is enabled. If an Activity Log Alert rule is not enabled, then none of its
	// actions will be activated.
	Enabled *bool
}

// WeeklyRecurrence - Weekly recurrence object.
type WeeklyRecurrence struct {
	// REQUIRED; Specifies the values for weekly recurrence pattern.
	DaysOfWeek []*DaysOfWeek

	// REQUIRED; Specifies when the recurrence should be applied.
	RecurrenceType *RecurrenceType

	// End time for recurrence.
	EndTime *string

	// Start time for recurrence.
	StartTime *string
}

// GetRecurrence implements the RecurrenceClassification interface for type WeeklyRecurrence.
func (w *WeeklyRecurrence) GetRecurrence() *Recurrence {
	return &Recurrence{
		EndTime:        w.EndTime,
		RecurrenceType: w.RecurrenceType,
		StartTime:      w.StartTime,
	}
}