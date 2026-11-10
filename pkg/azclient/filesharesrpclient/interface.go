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

// +azure:enableclientgen:=true
package filesharesrpclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares"
)

// Interface is the client for file shares of the Microsoft.FileShares resource provider
// (Microsoft.FileShares/fileShares). For file shares in a Microsoft.Storage storage account, use fileshareclient.
// Write methods wait until the long-running operation finishes.
// +azure:client:resource=FileSharesRP,packageName=github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares,packageAlias=armfileshares,clientName=Client,expand=false,crossSubFactory=true
type Interface interface {
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result *armfileshares.FileShare, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, resourceParam armfileshares.FileShare) (result *armfileshares.FileShare, err error)
	Update(ctx context.Context, resourceGroupName string, resourceName string, properties armfileshares.FileShareUpdate) (result *armfileshares.FileShare, err error)
	Delete(ctx context.Context, resourceGroupName string, resourceName string) (err error)
}
