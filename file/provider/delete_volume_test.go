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

// Package provider ...
package provider

import (
	"errors"
	"testing"

	"github.com/IBM/ibmcloud-volume-interface/lib/provider"
	util "github.com/IBM/ibmcloud-volume-interface/lib/utils"
	"github.com/IBM/ibmcloud-volume-interface/lib/utils/reasoncode"
	"github.com/IBM/ibmcloud-volume-vpc/common/vpcclient/models"
	serviceFakes "github.com/IBM/ibmcloud-volume-vpc/common/vpcclient/vpcvolume/fakes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestDeleteVolume(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		volumeService *serviceFakes.VolumeService
	)

	testCases := []struct {
		testCaseName   string
		baseVolume     *models.Volume
		providerVolume *provider.Volume

		tags  map[string]string
		setup func()

		skipErrTest        bool
		expectedErr        string
		expectedReasonCode string

		verify func(t *testing.T, err error)
	}{
		{
			testCaseName: "Not supported yet",
			providerVolume: &provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("Test volume"),
				Capacity: Int(10),
				Iops:     String("1000"),
				VPCVolume: provider.VPCVolume{
					Profile: &provider.Profile{Name: "general-purpose"},
				},
			},
			verify: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		}, {
			testCaseName:       "False positive: No volume being sent",
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description:'Not a valid volume ID",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Incorrect volume ID",
			providerVolume: &provider.Volume{
				VolumeID: "wrong volume ID",
				Name:     String("Test volume"),
				Capacity: Int(10),
				Iops:     String("1000"),
				VPCVolume: provider.VPCVolume{
					Profile:       &provider.Profile{Name: "general-purpose"},
					ResourceGroup: &provider.ResourceGroup{ID: "default resource group id", Name: "default resource group"},
				},
			},
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description:'Not a valid volume ID",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		}, {
			testCaseName: "Incorrect volume ID",
			providerVolume: &provider.Volume{
				VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
				Name:     String("Test volume"),
				Capacity: Int(10),
				Iops:     String("1000"),
				VPCVolume: provider.VPCVolume{
					Profile:       &provider.Profile{Name: "general-purpose"},
					ResourceGroup: &provider.ResourceGroup{ID: "default resource group id", Name: "default resource group"},
				},
			},
			expectedErr:        "{Code:ErrorUnclassified, Type:InvalidRequest, Description:'Not a valid volume ID",
			expectedReasonCode: "ErrorUnclassified",
			verify: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.testCaseName, func(t *testing.T) {
			vpcs, uc, sc, err := GetTestOpenSession(t, logger)
			assert.NotNil(t, vpcs)
			assert.NotNil(t, uc)
			assert.NotNil(t, sc)
			assert.Nil(t, err)

			volumeService = &serviceFakes.VolumeService{}
			assert.NotNil(t, volumeService)
			uc.VolumeServiceReturns(volumeService)

			if testcase.expectedErr != "" {
				volumeService.DeleteVolumeReturns(errors.New(testcase.expectedReasonCode))
				volumeService.GetVolumeReturns(testcase.baseVolume, errors.New(testcase.expectedReasonCode))
			} else {
				volumeService.DeleteVolumeReturns(nil)
				volumeService.GetVolumeReturns(testcase.baseVolume, nil)
			}

			err = vpcs.DeleteVolume(testcase.providerVolume)

			if testcase.expectedErr != "" {
				assert.NotNil(t, err)
				logger.Info("Error details", zap.Reflect("Error details", err.Error()))
				assert.Equal(t, reasoncode.ReasonCode(testcase.expectedReasonCode), util.ErrorReasonCode(err))
			}

			if testcase.verify != nil {
				testcase.verify(t, err)
			}
		})
	}
}

func TestDeleteVolumeTwo(t *testing.T) {
	//var err error
	logger, teardown := GetTestLogger(t)
	defer teardown()

	var (
		volumeService *serviceFakes.VolumeService
	)

	var baseVolume *models.Volume
	var providerVolume *provider.Volume

	baseVolume = &models.Volume{
		ID:     "16f293bf-test-4bff-816f-e199c0c65db5",
		Name:   "test volume name",
		Status: models.StatusType("OK"),
		Iops:   int64(1000),
		Zone:   &models.Zone{Name: "test-zone"},
	}

	providerVolume = &provider.Volume{
		VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
		Name:     String("Test volume"),
		Capacity: Int(10),
		Iops:     String("1000"),
		VPCVolume: provider.VPCVolume{
			Profile:       &provider.Profile{Name: "general-purpose"},
			ResourceGroup: &provider.ResourceGroup{ID: "default resource group id", Name: "default resource group"},
		},
	}

	vpcs, uc, sc, err := GetTestOpenSession(t, logger)
	assert.NotNil(t, vpcs)
	assert.NotNil(t, uc)
	assert.NotNil(t, sc)
	assert.Nil(t, err)

	volumeService = &serviceFakes.VolumeService{}
	assert.NotNil(t, volumeService)
	uc.VolumeServiceReturns(volumeService)

	volumeService.DeleteVolumeReturns(errors.New("not_found"))
	volumeService.GetVolumeReturns(nil, errors.New("not_found"))

	err = vpcs.DeleteVolume(providerVolume)
	assert.NotNil(t, err)

	volumeService.DeleteVolumeReturns(errors.New("failedToDeleteVolume"))
	volumeService.GetVolumeReturns(baseVolume, nil)

	err = vpcs.DeleteVolume(providerVolume)
	assert.NotNil(t, err)

	volumeService.DeleteVolumeReturns(errors.New("failedToDeleteVolume"))
	volumeService.GetVolumeReturns(nil, errors.New("wrong code"))

	err = vpcs.DeleteVolume(providerVolume)
	assert.NotNil(t, err)
}
