//
// Copyright (c) 2015-2024 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
//

package madmin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// ServerPeerUpdateStatus server update peer binary update result
type ServerPeerUpdateStatus struct {
	Host           string                 `json:"host"`
	Err            string                 `json:"err,omitempty"`
	CurrentVersion string                 `json:"currentVersion"`
	UpdatedVersion string                 `json:"updatedVersion"`
	WaitingDrives  map[string]DiskMetrics `json:"waitingDrives,omitempty"`
}

// ServerUpdateStatusV2 server update status
type ServerUpdateStatusV2 struct {
	DryRun  bool                     `json:"dryRun"`
	Results []ServerPeerUpdateStatus `json:"results,omitempty"`
}

// ServerUpdateV2 - updates and restarts the MinIO cluster to latest version.
// optionally takes an input URL to specify a custom update binary link
func (adm *AdminClient) ServerUpdateV2(ctx context.Context, updateURL string) (us ServerUpdateStatusV2, err error) {
	queryValues := url.Values{}
	queryValues.Set("updateURL", updateURL)
	queryValues.Set("type", "2")

	// Request API to Restart server
	resp, err := adm.executeMethod(ctx,
		http.MethodPost, requestData{
			relPath:     adminAPIPrefix + "/update",
			queryValues: queryValues,
		},
	)
	defer closeResponse(resp)
	if err != nil {
		return us, err
	}

	if resp.StatusCode != http.StatusOK {
		return us, httpRespToErrorResponse(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&us); err != nil {
		return us, err
	}

	return us, nil
}

// ServerUpdateStatus - contains the response of service update API
type ServerUpdateStatus struct {
	// Deprecated: this struct is fully deprecated since Jan 2024.
	CurrentVersion string `json:"currentVersion"`
	UpdatedVersion string `json:"updatedVersion"`
}

// ServerUpdate - updates and restarts the MinIO cluster to latest version.
// optionally takes an input URL to specify a custom update binary link
func (adm *AdminClient) ServerUpdate(ctx context.Context, updateURL string) (us ServerUpdateStatus, err error) {
	queryValues := url.Values{}
	queryValues.Set("updateURL", updateURL)

	// Request API to Restart server
	resp, err := adm.executeMethod(ctx,
		http.MethodPost, requestData{
			relPath:     adminAPIPrefix + "/update",
			queryValues: queryValues,
		},
	)
	defer closeResponse(resp)
	if err != nil {
		return us, err
	}

	if resp.StatusCode != http.StatusOK {
		return us, httpRespToErrorResponse(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&us); err != nil {
		return us, err
	}

	return us, nil
}
