package version

import "testing"

func TestVersionCompatibility(t *testing.T) {
	tests := []struct {
		name    string
		server  Version
		client  Version
		isValid bool
	}{
		{name: "same version", server: "v0.77.18.260814", client: "v0.77.18.260814", isValid: true},
		{name: "older build", server: "v0.77.18.260814", client: "v0.77.18.260813", isValid: true},
		{name: "newer build", server: "v0.77.18.260814", client: "v0.77.18.260815", isValid: false},
		{name: "older patch", server: "v0.77.18.260814", client: "v0.77.17.260814", isValid: false},
		{name: "newer patch", server: "v0.77.18.260814", client: "v0.77.19.260814", isValid: false},
		{name: "different minor", server: "v0.77.18.260814", client: "v0.78.18.260814", isValid: false},
		{name: "same three-part version", server: "v0.77.18.260814", client: "v0.77.18", isValid: true},
		{name: "different three-part version", server: "v0.77.18.260814", client: "v0.77.17", isValid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.Compatibility(tt.client); got != tt.isValid {
				t.Fatalf("Version(%q).Compatibility(%q) = %t, want %t", tt.server, tt.client, got, tt.isValid)
			}
		})
	}
}

func TestShortVersionsAndFormatting(t *testing.T) {
	for _, c := range []struct {
		version              Version
		tag, pkg, compatible string
	}{
		{"v1.2.3", "1.2.3", "1.2.3-0", "v1.2.3"},
		{"v1.2.3.260814-extra", "1.2.3-260814", "1.2.3-260814", "v1.2.3.260814"},
		{"1", "1.0.0", "1.0.0-0", "1.0.0"},
		{"1.2", "1.2.0", "1.2.0-0", "1.2.0"},
	} {
		if c.version.Tag() != c.tag || c.version.PkgSuffix() != c.pkg || c.version.CompatibleVersion().String() != c.compatible {
			t.Fatal(c, c.version.Tag(), c.version.PkgSuffix(), c.version.CompatibleVersion())
		}
	}
	for _, v := range []Version{"", "bad", "1..2", "1.2.3.1000000", "1000.0.0"} {
		if v.Compatibility(v) || v.Number() != 0 {
			t.Fatal("invalid version accepted", v)
		}
	}
}
