// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package version

import "github.com/hashicorp/packer-plugin-sdk/version"

var (
	Version           = "1.3.0"
	VersionPrerelease = "rc.5"
	VersionMetadata   = ""
	PluginVersion     = version.NewPluginVersion(Version, VersionPrerelease, VersionMetadata)
)
