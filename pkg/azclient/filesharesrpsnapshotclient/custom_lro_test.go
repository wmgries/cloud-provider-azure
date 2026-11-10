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
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares/fake"
	"github.com/stretchr/testify/require"
)

// failedOperation answers a write with one in-progress response, then a failed operation. The error exists only in
// the status poll, so a write that does not wait, or that drops the error, returns nil.
func failedOperation[T any]() (resp azfake.PollerResponder[T], errResp azfake.ErrorResponder) {
	resp.AddNonTerminalResponse(http.StatusAccepted, nil)
	resp.SetTerminalError(http.StatusConflict, "OperationFailed")
	return resp, errResp
}

func TestWritesReturnTheOperationError(t *testing.T) {
	client, err := New("00000000-0000-0000-0000-000000000000", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: fake.NewFileShareSnapshotsServerTransport(&fake.FileShareSnapshotsServer{
			BeginCreateOrUpdateFileShareSnapshot: func(context.Context, string, string, string, armfileshares.FileShareSnapshot, *armfileshares.FileShareSnapshotsClientBeginCreateOrUpdateFileShareSnapshotOptions) (azfake.PollerResponder[armfileshares.FileShareSnapshotsClientCreateOrUpdateFileShareSnapshotResponse], azfake.ErrorResponder) {
				return failedOperation[armfileshares.FileShareSnapshotsClientCreateOrUpdateFileShareSnapshotResponse]()
			},
			BeginDeleteFileShareSnapshot: func(context.Context, string, string, string, *armfileshares.FileShareSnapshotsClientBeginDeleteFileShareSnapshotOptions) (azfake.PollerResponder[armfileshares.FileShareSnapshotsClientDeleteFileShareSnapshotResponse], azfake.ErrorResponder) {
				return failedOperation[armfileshares.FileShareSnapshotsClientDeleteFileShareSnapshotResponse]()
			},
		})},
	})
	require.NoError(t, err)
	writes := map[string]func(context.Context) error{
		"CreateOrUpdate": func(ctx context.Context) error {
			_, err := client.CreateOrUpdate(ctx, "rg", "share", "snapshot", armfileshares.FileShareSnapshot{})
			return err
		},
		"Delete": func(ctx context.Context) error {
			return client.Delete(ctx, "rg", "share", "snapshot")
		},
	}
	for name, write := range writes {
		t.Run(name, func(t *testing.T) {
			var respErr *azcore.ResponseError
			require.ErrorAs(t, write(context.Background()), &respErr)
			require.Equal(t, "OperationFailed", respErr.ErrorCode)
		})
	}
}
