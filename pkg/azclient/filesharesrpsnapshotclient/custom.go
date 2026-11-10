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

package filesharesrpsnapshotclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/metrics"
	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/utils"
)

const (
	GetOperationName            = "FileSharesRPSnapshotClient.Get"
	CreateOrUpdateOperationName = "FileSharesRPSnapshotClient.CreateOrUpdate"
	DeleteOperationName         = "FileSharesRPSnapshotClient.Delete"
)

// Get gets a snapshot of a file share.
func (client *Client) Get(ctx context.Context, resourceGroupName string, parentResourceName string, resourceName string) (result *armfileshares.FileShareSnapshot, err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "FileSharesRPSnapshot", "get")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, GetOperationName, client.tracer, nil)
	defer endSpan(err)

	resp, err := client.GetFileShareSnapshot(ctx, resourceGroupName, parentResourceName, resourceName, nil)
	if err != nil {
		return nil, err
	}
	return &resp.FileShareSnapshot, nil
}

// CreateOrUpdate creates or updates a snapshot of a file share, and waits until the long-running operation finishes.
func (client *Client) CreateOrUpdate(ctx context.Context, resourceGroupName string, parentResourceName string, resourceName string, resourceParam armfileshares.FileShareSnapshot) (result *armfileshares.FileShareSnapshot, err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "FileSharesRPSnapshot", "create_or_update")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, CreateOrUpdateOperationName, client.tracer, nil)
	defer endSpan(err)

	resp, err := utils.NewPollerWrapper(client.BeginCreateOrUpdateFileShareSnapshot(ctx, resourceGroupName, parentResourceName, resourceName, resourceParam, nil)).WaitforPollerResp(ctx)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		return &resp.FileShareSnapshot, nil
	}
	return nil, nil
}

// Delete deletes a snapshot of a file share, and waits until the long-running operation finishes.
func (client *Client) Delete(ctx context.Context, resourceGroupName string, parentResourceName string, resourceName string) (err error) {
	metricsCtx := metrics.BeginARMRequest(client.subscriptionID, resourceGroupName, "FileSharesRPSnapshot", "delete")
	defer func() { metricsCtx.Observe(ctx, err) }()
	ctx, endSpan := runtime.StartSpan(ctx, DeleteOperationName, client.tracer, nil)
	defer endSpan(err)

	_, err = utils.NewPollerWrapper(client.BeginDeleteFileShareSnapshot(ctx, resourceGroupName, parentResourceName, resourceName, nil)).WaitforPollerResp(ctx)
	return err
}
