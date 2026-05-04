package installer

import "testing"

func TestIsNewerVersion(t *testing.T) {
	cases := []struct {
		name            string
		current, latest string
		want            bool
	}{
		// The bug that motivated this test: lexicographic compare reported
		// 1.0.10 as older than 1.0.9 because "1" < "9" as a string.
		{"patch double-digit beats single-digit", "1.0.9", "1.0.10", true},
		{"minor double-digit beats single-digit", "1.9.0", "1.10.0", true},
		{"major double-digit beats single-digit", "2.0.0", "10.0.0", true},

		{"equal versions are not newer", "1.2.3", "1.2.3", false},
		{"older patch is not newer", "1.2.3", "1.2.2", false},
		{"newer patch wins", "1.2.2", "1.2.3", true},

		{"v-prefix tolerated on both sides", "v1.0.0", "v1.0.1", true},
		{"v-prefix tolerated on one side", "1.0.0", "v1.0.1", true},

		// SemVer 2.0: a release outranks any pre-release of the same version.
		{"release beats pre-release of same version", "1.0.0-rc1", "1.0.0", true},
		{"pre-release loses to release of same version", "1.0.0", "1.0.0-rc1", false},
		{"newer pre-release wins", "1.0.0-rc1", "1.0.0-rc2", true},

		// Dev / dirty builds suppress the check entirely.
		{"dev build never has updates", "dev", "1.0.0", false},
		{"dirty build never has updates", "1.0.0-dirty", "1.0.0", false},
		{"empty current never has updates", "", "1.0.0", false},
	}

	v := &VersionChecker{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := v.isNewerVersion(tc.latest, tc.current)
			if got != tc.want {
				t.Errorf("isNewerVersion(latest=%q, current=%q) = %v, want %v",
					tc.latest, tc.current, got, tc.want)
			}
		})
	}
}
