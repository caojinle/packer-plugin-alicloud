// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package ecs

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/template/config"
)

func TestStepRegionCopyAlicloudImage_buildCopyImageRequest_sameRegionKMSKeyId(t *testing.T) {
	step := &stepRegionCopyAlicloudImage{
		RegionId:                      "cn-hangzhou",
		AlicloudImageDestinationNames: []string{"test-image"},
	}
	cfg := &Config{}
	cfg.AlicloudImageName = "test-image"
	cfg.ImageEncrypted = config.TrileanFromBool(true)
	cfg.KMSKeyId = "same-region-kms-key"

	req := step.buildCopyImageRequest("cn-hangzhou", cfg, "m-source", -1)

	if req.DestinationRegionId != "cn-hangzhou" {
		t.Fatalf("expected destination region cn-hangzhou, got %s", req.DestinationRegionId)
	}
	if req.DestinationImageName != "test-image" {
		t.Fatalf("expected destination image name test-image, got %s", req.DestinationImageName)
	}
	if !req.Encrypted.HasValue() {
		t.Fatal("expected encrypted to be set")
	}
	encrypted, err := req.Encrypted.GetValue()
	if err != nil || !encrypted {
		t.Fatalf("expected encrypted to be true, got %v, err: %v", req.Encrypted, err)
	}
	if req.KMSKeyId != "same-region-kms-key" {
		t.Fatalf("expected KMS key id same-region-kms-key, got %s", req.KMSKeyId)
	}
}

func TestStepRegionCopyAlicloudImage_buildCopyImageRequest_sameRegionDefaultKMS(t *testing.T) {
	step := &stepRegionCopyAlicloudImage{RegionId: "cn-hangzhou"}
	cfg := &Config{}
	cfg.AlicloudImageName = "test-image"
	cfg.ImageEncrypted = config.TrileanFromBool(true)

	req := step.buildCopyImageRequest("cn-hangzhou", cfg, "m-source", -1)

	if !req.Encrypted.HasValue() {
		t.Fatal("expected encrypted to be set")
	}
	encrypted, err := req.Encrypted.GetValue()
	if err != nil || !encrypted {
		t.Fatalf("expected encrypted to be true, got %v, err: %v", req.Encrypted, err)
	}
	if req.KMSKeyId != "" {
		t.Fatalf("expected default KMS key id, got %s", req.KMSKeyId)
	}
}

func TestStepRegionCopyAlicloudImage_buildCopyImageRequest_crossRegionKMSKeyIds(t *testing.T) {
	step := &stepRegionCopyAlicloudImage{
		RegionId:                      "cn-hangzhou",
		AlicloudImageDestinationNames: []string{"test-image-sh", "test-image-bj"},
	}
	cfg := &Config{}
	cfg.AlicloudImageName = "test-image"
	cfg.ImageEncrypted = config.TrileanFromBool(true)
	cfg.ImageCopyKMSKeyIds = []string{"key-sh", "key-bj"}

	req0 := step.buildCopyImageRequest("cn-shanghai", cfg, "m-source", 0)
	if req0.DestinationImageName != "test-image-sh" {
		t.Fatalf("expected destination image name test-image-sh, got %s", req0.DestinationImageName)
	}
	if req0.KMSKeyId != "key-sh" {
		t.Fatalf("expected KMS key id key-sh, got %s", req0.KMSKeyId)
	}

	req1 := step.buildCopyImageRequest("cn-beijing", cfg, "m-source", 1)
	if req1.DestinationImageName != "test-image-bj" {
		t.Fatalf("expected destination image name test-image-bj, got %s", req1.DestinationImageName)
	}
	if req1.KMSKeyId != "key-bj" {
		t.Fatalf("expected KMS key id key-bj, got %s", req1.KMSKeyId)
	}
}

func TestStepRegionCopyAlicloudImage_buildCopyImageRequest_crossRegionKMSKeyIdsShorterThanRegions(t *testing.T) {
	step := &stepRegionCopyAlicloudImage{RegionId: "cn-hangzhou"}
	cfg := &Config{}
	cfg.AlicloudImageName = "test-image"
	cfg.ImageEncrypted = config.TrileanFromBool(true)
	cfg.AlicloudImageDestinationRegions = []string{"cn-shanghai", "cn-beijing"}
	cfg.ImageCopyKMSKeyIds = []string{"key-sh"}

	req0 := step.buildCopyImageRequest("cn-shanghai", cfg, "m-source", 0)
	if req0.KMSKeyId != "key-sh" {
		t.Fatalf("expected KMS key id key-sh, got %s", req0.KMSKeyId)
	}

	req1 := step.buildCopyImageRequest("cn-beijing", cfg, "m-source", 1)
	if req1.KMSKeyId != "" {
		t.Fatalf("expected default KMS key id for missing entry, got %s", req1.KMSKeyId)
	}
}

func TestStepRegionCopyAlicloudImage_buildCopyImageRequest_crossRegionEmptyKMSKeyIdEntry(t *testing.T) {
	step := &stepRegionCopyAlicloudImage{RegionId: "cn-hangzhou"}
	cfg := &Config{}
	cfg.AlicloudImageName = "test-image"
	cfg.ImageEncrypted = config.TrileanFromBool(true)
	cfg.AlicloudImageDestinationRegions = []string{"cn-shanghai", "cn-beijing"}
	cfg.ImageCopyKMSKeyIds = []string{"", "key-bj"}

	req0 := step.buildCopyImageRequest("cn-shanghai", cfg, "m-source", 0)
	if req0.KMSKeyId != "" {
		t.Fatalf("expected default KMS key id for empty entry, got %s", req0.KMSKeyId)
	}

	req1 := step.buildCopyImageRequest("cn-beijing", cfg, "m-source", 1)
	if req1.KMSKeyId != "key-bj" {
		t.Fatalf("expected KMS key id key-bj, got %s", req1.KMSKeyId)
	}
}

func TestStepRegionCopyAlicloudImage_buildCopyImageRequest_noEncryption(t *testing.T) {
	step := &stepRegionCopyAlicloudImage{RegionId: "cn-hangzhou"}
	cfg := &Config{}
	cfg.AlicloudImageName = "test-image"

	req := step.buildCopyImageRequest("cn-hangzhou", cfg, "m-source", -1)

	if req.Encrypted.HasValue() {
		t.Fatalf("expected encrypted not to be set, got %v", req.Encrypted)
	}
	if req.KMSKeyId != "" {
		t.Fatalf("expected empty KMS key id, got %s", req.KMSKeyId)
	}
}

func TestStepRegionCopyAlicloudImage_buildCopyImageRequest_namesShorterThanRegionsWithEncryption(t *testing.T) {
	// Reproduces the template where image_copy_names has fewer entries than
	// image_copy_regions while image_encrypted is enabled: the encrypted
	// same-region copy must still take over image_name, and the cross-region
	// destination without a matching name entry must be left unnamed.
	step := &stepRegionCopyAlicloudImage{
		RegionId:                      "cn-hangzhou",
		AlicloudImageDestinationNames: []string{"custom-78d1328a"},
	}
	cfg := &Config{}
	cfg.AlicloudImageName = "39f339bf"
	cfg.ImageEncrypted = config.TrileanFromBool(true)

	req0 := step.buildCopyImageRequest("cn-shanghai", cfg, "m-source", 0)
	if req0.DestinationImageName != "custom-78d1328a" {
		t.Fatalf("expected destination image name custom-78d1328a, got %s", req0.DestinationImageName)
	}

	req1 := step.buildCopyImageRequest("cn-beijing", cfg, "m-source", 1)
	if req1.DestinationImageName != "" {
		t.Fatalf("expected empty destination image name for cn-beijing, got %s", req1.DestinationImageName)
	}

	req2 := step.buildCopyImageRequest("cn-hangzhou", cfg, "m-source", 1)
	if req2.DestinationImageName != "39f339bf" {
		t.Fatalf("expected destination image name 39f339bf, got %s", req2.DestinationImageName)
	}
}
