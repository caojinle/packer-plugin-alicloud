// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ecs

import (
	"context"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepImageDeleteSSHPrivateKey struct {
	AlicloudImageDeleteSSHPrivateKey bool
}

func (s *stepImageDeleteSSHPrivateKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)
	client := state.Get("client").(*ClientWrapper)
	config := state.Get("config").(*Config)
	instance := state.Get("instance").(*ecs.Instance)
	keyPairName := config.Comm.SSHKeyPairName

	// Check for force delete
	if s.AlicloudImageDeleteSSHPrivateKey {
		if keyPairName == "" {
			return multistep.ActionContinue
		}
		detachKeyPairRequest := ecs.CreateDetachKeyPairRequest()
		detachKeyPairRequest.RegionId = config.AlicloudRegion
		detachKeyPairRequest.KeyPairName = keyPairName
		detachKeyPairRequest.InstanceIds = fmt.Sprintf("[\"%s\"]", instance.InstanceId)

		if _, err := client.DetachKeyPair(detachKeyPairRequest); err != nil {
			return halt(state, err, fmt.Sprintf("Error Detaching keypair %s to instance %s", keyPairName, instance.InstanceId))
		}
		ui.Say(fmt.Sprintf("Detach keypair %s from instance: %s", keyPairName, instance.InstanceId))

		rebootInstanceRequest := ecs.CreateRebootInstanceRequest()
		rebootInstanceRequest.InstanceId = instance.InstanceId
		if _, err := client.RebootInstance(rebootInstanceRequest); err != nil {
			return halt(state, err, "Error reboot instance")
		}

		_, err := client.WaitForInstanceStatus(instance.RegionId, instance.InstanceId, InstanceStatusRunning)
		if err != nil {
			return halt(state, err, "Timeout waiting for instance to reboot")
		}
		ui.Say(fmt.Sprintf("Reboot instance: %s", instance.InstanceId))

	}
	return multistep.ActionContinue
}

func (s *stepImageDeleteSSHPrivateKey) Cleanup(multistep.StateBag) {
	// No cleanup...
}
