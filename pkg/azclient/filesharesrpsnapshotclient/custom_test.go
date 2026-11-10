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
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"sigs.k8s.io/cloud-provider-azure/pkg/azclient/utils"
)

var shareClient *armfileshares.Client

func init() {
	// Microsoft.FileShares does not run in Azure China, so there is no AZURECHINACLOUD cassette.
	// A suite with no specs does not start the recorder.
	if strings.EqualFold(os.Getenv("AZURE_CLOUD"), "AZURECHINACLOUD") {
		return
	}
	additionalTestCases = func() {
		ginkgo.When("creation requests are raised", func() {
			ginkgo.It("should return the snapshot after the operation finishes", func(ctx context.Context) {
				resourceName = "akscitsnapshot"
				snapshot, err := realClient.CreateOrUpdate(ctx, resourceGroupName, fileshareName, resourceName, armfileshares.FileShareSnapshot{
					Properties: &armfileshares.FileShareSnapshotProperties{Metadata: map[string]*string{"test": to.Ptr("snapshot")}},
				})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(snapshot).To(gomega.HaveField("Name", gomega.HaveValue(gomega.Equal(resourceName))))
			})
		})
		ginkgo.When("get requests are raised", func() {
			ginkgo.It("should return the snapshot", func(ctx context.Context) {
				snapshot, err := realClient.Get(ctx, resourceGroupName, fileshareName, resourceName)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(snapshot).To(gomega.HaveField("Name", gomega.HaveValue(gomega.Equal(resourceName))))
			})
		})
		ginkgo.When("deletion requests are raised", func() {
			ginkgo.It("should delete the snapshot, so that a later get returns 404", func(ctx context.Context) {
				gomega.Expect(realClient.Delete(ctx, resourceGroupName, fileshareName, resourceName)).To(gomega.Succeed())
				snapshot, err := realClient.Get(ctx, resourceGroupName, fileshareName, resourceName)
				var respErr *azcore.ResponseError
				gomega.Expect(errors.As(err, &respErr)).To(gomega.BeTrue(), "got %v", err)
				gomega.Expect(respErr.StatusCode).To(gomega.Equal(http.StatusNotFound))
				gomega.Expect(snapshot).To(gomega.BeNil())
			})
		})
	}

	beforeAllFunc = func(ctx context.Context) {
		fileshareName = "akscitsnapshotshare"
		shareClient, err = armfileshares.NewClient(subscriptionID, recorder.TokenCredential(), &arm.ClientOptions{ClientOptions: clientOption})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		_, err = utils.NewPollerWrapper(shareClient.BeginCreateOrUpdate(ctx, resourceGroupName, fileshareName, armfileshares.FileShare{
			Location: to.Ptr(location),
			Properties: &armfileshares.FileShareProperties{
				MountName:                      to.Ptr("fileshare"),
				MediaTier:                      to.Ptr(armfileshares.MediaTierSSD),
				Redundancy:                     to.Ptr(armfileshares.RedundancyLocal),
				Protocol:                       to.Ptr(armfileshares.ProtocolNFS),
				ProvisionedStorageGiB:          to.Ptr[int32](32),
				ProvisionedIOPerSec:            to.Ptr[int32](3000),
				ProvisionedThroughputMiBPerSec: to.Ptr[int32](125),
				NfsProtocolProperties:          &armfileshares.NfsProtocolProperties{RootSquash: to.Ptr(armfileshares.ShareRootSquashNoRootSquash)},
				PublicNetworkAccess:            to.Ptr(armfileshares.PublicNetworkAccessEnabled),
			},
		}, nil)).WaitforPollerResp(ctx)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}
	afterAllFunc = func(ctx context.Context) {
		_, err = utils.NewPollerWrapper(shareClient.BeginDelete(ctx, resourceGroupName, fileshareName, nil)).WaitforPollerResp(ctx)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}
}
