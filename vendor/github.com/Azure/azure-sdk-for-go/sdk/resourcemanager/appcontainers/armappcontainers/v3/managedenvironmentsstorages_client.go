//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armappcontainers

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ManagedEnvironmentsStoragesClient contains the methods for the ManagedEnvironmentsStorages group.
// Don't use this type directly, use NewManagedEnvironmentsStoragesClient() instead.
type ManagedEnvironmentsStoragesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewManagedEnvironmentsStoragesClient creates a new instance of ManagedEnvironmentsStoragesClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewManagedEnvironmentsStoragesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ManagedEnvironmentsStoragesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ManagedEnvironmentsStoragesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update storage for a managedEnvironment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - environmentName - Name of the Environment.
//   - storageName - Name of the storage.
//   - storageEnvelope - Configuration details of storage.
//   - options - ManagedEnvironmentsStoragesClientCreateOrUpdateOptions contains the optional parameters for the ManagedEnvironmentsStoragesClient.CreateOrUpdate
//     method.
func (client *ManagedEnvironmentsStoragesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, environmentName string, storageName string, storageEnvelope ManagedEnvironmentStorage, options *ManagedEnvironmentsStoragesClientCreateOrUpdateOptions) (ManagedEnvironmentsStoragesClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "ManagedEnvironmentsStoragesClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, environmentName, storageName, storageEnvelope, options)
	if err != nil {
		return ManagedEnvironmentsStoragesClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ManagedEnvironmentsStoragesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ManagedEnvironmentsStoragesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ManagedEnvironmentsStoragesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, environmentName string, storageName string, storageEnvelope ManagedEnvironmentStorage, options *ManagedEnvironmentsStoragesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/managedEnvironments/{environmentName}/storages/{storageName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if environmentName == "" {
		return nil, errors.New("parameter environmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{environmentName}", url.PathEscape(environmentName))
	if storageName == "" {
		return nil, errors.New("parameter storageName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageName}", url.PathEscape(storageName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, storageEnvelope); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ManagedEnvironmentsStoragesClient) createOrUpdateHandleResponse(resp *http.Response) (ManagedEnvironmentsStoragesClientCreateOrUpdateResponse, error) {
	result := ManagedEnvironmentsStoragesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedEnvironmentStorage); err != nil {
		return ManagedEnvironmentsStoragesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete storage for a managedEnvironment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - environmentName - Name of the Environment.
//   - storageName - Name of the storage.
//   - options - ManagedEnvironmentsStoragesClientDeleteOptions contains the optional parameters for the ManagedEnvironmentsStoragesClient.Delete
//     method.
func (client *ManagedEnvironmentsStoragesClient) Delete(ctx context.Context, resourceGroupName string, environmentName string, storageName string, options *ManagedEnvironmentsStoragesClientDeleteOptions) (ManagedEnvironmentsStoragesClientDeleteResponse, error) {
	var err error
	const operationName = "ManagedEnvironmentsStoragesClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, environmentName, storageName, options)
	if err != nil {
		return ManagedEnvironmentsStoragesClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ManagedEnvironmentsStoragesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return ManagedEnvironmentsStoragesClientDeleteResponse{}, err
	}
	return ManagedEnvironmentsStoragesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ManagedEnvironmentsStoragesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, environmentName string, storageName string, options *ManagedEnvironmentsStoragesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/managedEnvironments/{environmentName}/storages/{storageName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if environmentName == "" {
		return nil, errors.New("parameter environmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{environmentName}", url.PathEscape(environmentName))
	if storageName == "" {
		return nil, errors.New("parameter storageName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageName}", url.PathEscape(storageName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get storage for a managedEnvironment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - environmentName - Name of the Environment.
//   - storageName - Name of the storage.
//   - options - ManagedEnvironmentsStoragesClientGetOptions contains the optional parameters for the ManagedEnvironmentsStoragesClient.Get
//     method.
func (client *ManagedEnvironmentsStoragesClient) Get(ctx context.Context, resourceGroupName string, environmentName string, storageName string, options *ManagedEnvironmentsStoragesClientGetOptions) (ManagedEnvironmentsStoragesClientGetResponse, error) {
	var err error
	const operationName = "ManagedEnvironmentsStoragesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, environmentName, storageName, options)
	if err != nil {
		return ManagedEnvironmentsStoragesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ManagedEnvironmentsStoragesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ManagedEnvironmentsStoragesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ManagedEnvironmentsStoragesClient) getCreateRequest(ctx context.Context, resourceGroupName string, environmentName string, storageName string, options *ManagedEnvironmentsStoragesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/managedEnvironments/{environmentName}/storages/{storageName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if environmentName == "" {
		return nil, errors.New("parameter environmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{environmentName}", url.PathEscape(environmentName))
	if storageName == "" {
		return nil, errors.New("parameter storageName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageName}", url.PathEscape(storageName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ManagedEnvironmentsStoragesClient) getHandleResponse(resp *http.Response) (ManagedEnvironmentsStoragesClientGetResponse, error) {
	result := ManagedEnvironmentsStoragesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedEnvironmentStorage); err != nil {
		return ManagedEnvironmentsStoragesClientGetResponse{}, err
	}
	return result, nil
}

// List - Get all storages for a managedEnvironment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - environmentName - Name of the Environment.
//   - options - ManagedEnvironmentsStoragesClientListOptions contains the optional parameters for the ManagedEnvironmentsStoragesClient.List
//     method.
func (client *ManagedEnvironmentsStoragesClient) List(ctx context.Context, resourceGroupName string, environmentName string, options *ManagedEnvironmentsStoragesClientListOptions) (ManagedEnvironmentsStoragesClientListResponse, error) {
	var err error
	const operationName = "ManagedEnvironmentsStoragesClient.List"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.listCreateRequest(ctx, resourceGroupName, environmentName, options)
	if err != nil {
		return ManagedEnvironmentsStoragesClientListResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ManagedEnvironmentsStoragesClientListResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ManagedEnvironmentsStoragesClientListResponse{}, err
	}
	resp, err := client.listHandleResponse(httpResp)
	return resp, err
}

// listCreateRequest creates the List request.
func (client *ManagedEnvironmentsStoragesClient) listCreateRequest(ctx context.Context, resourceGroupName string, environmentName string, options *ManagedEnvironmentsStoragesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/managedEnvironments/{environmentName}/storages"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if environmentName == "" {
		return nil, errors.New("parameter environmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{environmentName}", url.PathEscape(environmentName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ManagedEnvironmentsStoragesClient) listHandleResponse(resp *http.Response) (ManagedEnvironmentsStoragesClientListResponse, error) {
	result := ManagedEnvironmentsStoragesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedEnvironmentStoragesCollection); err != nil {
		return ManagedEnvironmentsStoragesClientListResponse{}, err
	}
	return result, nil
}