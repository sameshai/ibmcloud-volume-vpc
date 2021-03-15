/**
 * Copyright 2020 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package vpcvolume_test ...
package vpcfilevolume_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/IBM/ibmcloud-volume-vpc/common/vpcclient/models"
	"github.com/IBM/ibmcloud-volume-vpc/common/vpcclient/riaas/test"
	"github.com/IBM/ibmcloud-volume-vpc/common/vpcclient/vpcfilevolume"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestListVolumes(t *testing.T) {
	// Setup new style zap logger
	logger, _ := GetTestContextLogger()
	defer logger.Sync()

	testCases := []struct {
		name string

		// Response
		status  int
		content string

		limit   int
		start   string
		filters *models.ListVolumeFilters

		// Expected return
		expectErr string
		verify    func(*testing.T)
		muxVerify func(*testing.T, *http.Request)
	}{
		{
			name:   "Verify that the correct endpoint is invoked",
			status: http.StatusNoContent,
		}, {
			name:      "Verify that a 404 is returned to the caller",
			status:    http.StatusNotFound,
			content:   "{\"errors\":[{\"message\":\"testerr\"}]}",
			expectErr: "Trace Code:, testerr Please check ",
		}, {
			name:   "Verify that limit is added to the query",
			limit:  12,
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"limit": []string{"12"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name:   "Verify that start is added to the query",
			start:  "x-y-z",
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"start": []string{"x-y-z"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name: "Verify that resource_group.id is added to the query",
			filters: &models.ListVolumeFilters{
				ResourceGroupID: "rgid",
			},
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"resource_group.id": []string{"rgid"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name: "Verify that tag is added to the query",
			filters: &models.ListVolumeFilters{
				Tag: "taggg",
			},
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"tag": []string{"taggg"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name: "Verify that zone.name is added to the query",
			filters: &models.ListVolumeFilters{
				ZoneName: "zone",
			},
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"zone.name": []string{"zone"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		}, {
			name: "Verify that volume name is added to the query",
			filters: &models.ListVolumeFilters{
				VolumeName: "testname",
			},
			status: http.StatusNoContent,
			muxVerify: func(t *testing.T, r *http.Request) {
				expectedValues := url.Values{"name": []string{"testname"}, "version": []string{models.APIVersion}}
				actualValues := r.URL.Query()
				assert.Equal(t, expectedValues, actualValues)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			mux, client, teardown := test.SetupServer(t)
			test.SetupMuxResponse(t, mux, vpcfilevolume.Version+"/shares", http.MethodGet, nil, testcase.status, testcase.content, nil)

			defer teardown()

			logger.Info("Test case being executed", zap.Reflect("testcase", testcase.name))

			volumeFileService := vpcfilevolume.New(client)

			volumes, err := volumeFileService.ListVolumes(testcase.limit, testcase.start, testcase.filters, logger)
			logger.Info("Volumes", zap.Reflect("volumes", volumes))

			if testcase.expectErr != "" && assert.Error(t, err) {
				assert.Equal(t, testcase.expectErr, err.Error())
				assert.Nil(t, volumes)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, volumes)
			}

			if testcase.verify != nil {
				testcase.verify(t)
			}
		})
	}
}
