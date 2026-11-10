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
package filesharesrpsnapshotclient

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares"
)

// Interface is the client for snapshots of Microsoft.FileShares file shares
// (Microsoft.FileShares/fileShares/fileShareSnapshots). parentResourceName is the file share name.
// Write methods wait until the long-running operation finishes.
// +azure:client:resource=FileShare,subResource=FileSharesRPSnapshot,packageName=github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares,packageAlias=armfileshares,clientName=FileShareSnapshotsClient,expand=false,crossSubFactory=true
type Interface interface {
	Get(ctx context.Context, resourceGroupName string, parentResourceName string, resourceName string) (result *armfileshares.FileShareSnapshot, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, parentResourceName string, resourceName string, resourceParam armfileshares.FileShareSnapshot) (result *armfileshares.FileShareSnapshot, err error)
	Delete(ctx context.Context, resourceGroupName string, parentResourceName string, resourceName string) (err error)
}
