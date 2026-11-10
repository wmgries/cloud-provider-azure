/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package filesharesrpclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/metrics"
	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/utils"
)

const (
	GetOperationName            = "FileSharesRPClient.Get"
	CreateOrUpdateOperationName = "FileSharesRPClient.CreateOrUpdate"
	UpdateOperationName         = "FileSharesRPClient.Update"
	DeleteOperationName         = "FileSharesRPClient.Delete"
)

// Get gets a file share.
func (client *Client) Get(ctx context.Context, resourceGroupName string, resourceName string) (result *armfileshares.FileShare, err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "FileSharesRP", "get")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, GetOperationName, client.tracer, nil)
	defer endSpan(err)

	resp, err := client.Client.Get(ctx, resourceGroupName, resourceName, nil)
	if err != nil {
		return nil, err
	}
	return &resp.FileShare, nil
}

// CreateOrUpdate creates or updates a file share, and waits until the long-running operation finishes.
func (client *Client) CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, resourceParam armfileshares.FileShare) (result *armfileshares.FileShare, err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "FileSharesRP", "create_or_update")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, CreateOrUpdateOperationName, client.tracer, nil)
	defer endSpan(err)

	resp, err := utils.NewPollerWrapper(client.Client.BeginCreateOrUpdate(ctx, resourceGroupName, resourceName, resourceParam, nil)).WaitforPollerResp(ctx)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		return &resp.FileShare, nil
	}
	return nil, nil
}

// Update changes the properties of a file share with a partial update, and waits until the long-running
// operation finishes.
func (client *Client) Update(ctx context.Context, resourceGroupName string, resourceName string, properties armfileshares.FileShareUpdate) (result *armfileshares.FileShare, err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "FileSharesRP", "update")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, UpdateOperationName, client.tracer, nil)
	defer endSpan(err)

	resp, err := utils.NewPollerWrapper(client.Client.BeginUpdate(ctx, resourceGroupName, resourceName, properties, nil)).WaitforPollerResp(ctx)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		return &resp.FileShare, nil
	}
	return nil, nil
}

// Delete deletes a file share, and waits until the long-running operation finishes.
func (client *Client) Delete(ctx context.Context, resourceGroupName string, resourceName string) (err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "FileSharesRP", "delete")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, DeleteOperationName, client.tracer, nil)
	defer endSpan(err)

	_, err = utils.NewPollerWrapper(client.Client.BeginDelete(ctx, resourceGroupName, resourceName, nil)).WaitforPollerResp(ctx)
	return err
}
