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
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func init() {
	// Microsoft.FileShares does not run in Azure China, so there is no AZURECHINACLOUD cassette.
	// A suite with no specs does not start the recorder.
	if strings.EqualFold(os.Getenv("AZURE_CLOUD"), "AZURECHINACLOUD") {
		return
	}
	additionalTestCases = func() {
		ginkgo.When("creation requests are raised", func() {
			ginkgo.It("should return the file share after the operation finishes", func(ctx context.Context) {
				resourceName = "akscitfilesharesrp"
				share, err := realClient.CreateOrUpdate(ctx, resourceGroupName, resourceName, armfileshares.FileShare{
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
				})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(share).To(gomega.HaveField("Name", gomega.HaveValue(gomega.Equal(resourceName))))
				gomega.Expect(share).To(gomega.HaveField("Properties.ProvisioningState", gomega.HaveValue(gomega.Equal(armfileshares.FileShareProvisioningStateSucceeded))))
			})
		})
		ginkgo.When("get requests are raised", func() {
			ginkgo.It("should return the file share", func(ctx context.Context) {
				share, err := realClient.Get(ctx, resourceGroupName, resourceName)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(share).To(gomega.HaveField("Name", gomega.HaveValue(gomega.Equal(resourceName))))
			})
		})
		ginkgo.When("update requests are raised", func() {
			ginkgo.It("should return the file share with the new IOPS after the operation finishes", func(ctx context.Context) {
				share, err := realClient.Update(ctx, resourceGroupName, resourceName, armfileshares.FileShareUpdate{
					Properties: &armfileshares.FileShareUpdateProperties{ProvisionedIOPerSec: to.Ptr[int32](3200)},
				})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(share).To(gomega.HaveField("Properties.ProvisionedIOPerSec", gomega.HaveValue(gomega.Equal(int32(3200)))))
			})
		})
		ginkgo.When("deletion requests are raised", func() {
			ginkgo.It("should delete the file share, so that a later get returns 404", func(ctx context.Context) {
				gomega.Expect(realClient.Delete(ctx, resourceGroupName, resourceName)).To(gomega.Succeed())
				share, err := realClient.Get(ctx, resourceGroupName, resourceName)
				var respErr *azcore.ResponseError
				gomega.Expect(errors.As(err, &respErr)).To(gomega.BeTrue(), "got %v", err)
				gomega.Expect(respErr.StatusCode).To(gomega.Equal(http.StatusNotFound))
				gomega.Expect(share).To(gomega.BeNil())
			})
		})
	}
}
