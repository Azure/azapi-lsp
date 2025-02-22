//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armrecoveryservicessiterecovery

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ReplicationVaultSettingClient contains the methods for the ReplicationVaultSetting group.
// Don't use this type directly, use NewReplicationVaultSettingClient() instead.
type ReplicationVaultSettingClient struct {
	host              string
	resourceName      string
	resourceGroupName string
	subscriptionID    string
	pl                runtime.Pipeline
}

// NewReplicationVaultSettingClient creates a new instance of ReplicationVaultSettingClient with the specified values.
// resourceName - The name of the recovery services vault.
// resourceGroupName - The name of the resource group where the recovery services vault is present.
// subscriptionID - The subscription Id.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewReplicationVaultSettingClient(resourceName string, resourceGroupName string, subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ReplicationVaultSettingClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &ReplicationVaultSettingClient{
		resourceName:      resourceName,
		resourceGroupName: resourceGroupName,
		subscriptionID:    subscriptionID,
		host:              ep,
		pl:                pl,
	}
	return client, nil
}

// BeginCreate - The operation to configure vault setting.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-10-01
// vaultSettingName - Vault setting name.
// input - Vault setting creation input.
// options - ReplicationVaultSettingClientBeginCreateOptions contains the optional parameters for the ReplicationVaultSettingClient.BeginCreate
// method.
func (client *ReplicationVaultSettingClient) BeginCreate(ctx context.Context, vaultSettingName string, input VaultSettingCreationInput, options *ReplicationVaultSettingClientBeginCreateOptions) (*runtime.Poller[ReplicationVaultSettingClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, vaultSettingName, input, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[ReplicationVaultSettingClientCreateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[ReplicationVaultSettingClientCreateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Create - The operation to configure vault setting.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-10-01
func (client *ReplicationVaultSettingClient) create(ctx context.Context, vaultSettingName string, input VaultSettingCreationInput, options *ReplicationVaultSettingClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, vaultSettingName, input, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *ReplicationVaultSettingClient) createCreateRequest(ctx context.Context, vaultSettingName string, input VaultSettingCreationInput, options *ReplicationVaultSettingClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationVaultSettings/{vaultSettingName}"
	if client.resourceName == "" {
		return nil, errors.New("parameter client.resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(client.resourceName))
	if client.resourceGroupName == "" {
		return nil, errors.New("parameter client.resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(client.resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if vaultSettingName == "" {
		return nil, errors.New("parameter vaultSettingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultSettingName}", url.PathEscape(vaultSettingName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, input)
}

// Get - Gets the vault setting. This includes the Migration Hub connection settings.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-10-01
// vaultSettingName - Vault setting name.
// options - ReplicationVaultSettingClientGetOptions contains the optional parameters for the ReplicationVaultSettingClient.Get
// method.
func (client *ReplicationVaultSettingClient) Get(ctx context.Context, vaultSettingName string, options *ReplicationVaultSettingClientGetOptions) (ReplicationVaultSettingClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, vaultSettingName, options)
	if err != nil {
		return ReplicationVaultSettingClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ReplicationVaultSettingClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ReplicationVaultSettingClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ReplicationVaultSettingClient) getCreateRequest(ctx context.Context, vaultSettingName string, options *ReplicationVaultSettingClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationVaultSettings/{vaultSettingName}"
	if client.resourceName == "" {
		return nil, errors.New("parameter client.resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(client.resourceName))
	if client.resourceGroupName == "" {
		return nil, errors.New("parameter client.resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(client.resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if vaultSettingName == "" {
		return nil, errors.New("parameter vaultSettingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultSettingName}", url.PathEscape(vaultSettingName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ReplicationVaultSettingClient) getHandleResponse(resp *http.Response) (ReplicationVaultSettingClientGetResponse, error) {
	result := ReplicationVaultSettingClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VaultSetting); err != nil {
		return ReplicationVaultSettingClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets the list of vault setting. This includes the Migration Hub connection settings.
// Generated from API version 2022-10-01
// options - ReplicationVaultSettingClientListOptions contains the optional parameters for the ReplicationVaultSettingClient.List
// method.
func (client *ReplicationVaultSettingClient) NewListPager(options *ReplicationVaultSettingClientListOptions) *runtime.Pager[ReplicationVaultSettingClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReplicationVaultSettingClientListResponse]{
		More: func(page ReplicationVaultSettingClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReplicationVaultSettingClientListResponse) (ReplicationVaultSettingClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ReplicationVaultSettingClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ReplicationVaultSettingClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ReplicationVaultSettingClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ReplicationVaultSettingClient) listCreateRequest(ctx context.Context, options *ReplicationVaultSettingClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationVaultSettings"
	if client.resourceName == "" {
		return nil, errors.New("parameter client.resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(client.resourceName))
	if client.resourceGroupName == "" {
		return nil, errors.New("parameter client.resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(client.resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ReplicationVaultSettingClient) listHandleResponse(resp *http.Response) (ReplicationVaultSettingClientListResponse, error) {
	result := ReplicationVaultSettingClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VaultSettingCollection); err != nil {
		return ReplicationVaultSettingClientListResponse{}, err
	}
	return result, nil
}
