package tfsec

import (
	"testing"

	"github.com/k1rd3rf/tfsec/internal/app/tfsec/scanner"

	"github.com/k1rd3rf/tfsec/internal/app/tfsec/checks"
)

func Test_GkeAbacEnabled(t *testing.T) {

	var tests = []struct {
		name                  string
		source                string
		mustIncludeResultCode scanner.RuleID
		mustExcludeResultCode scanner.RuleID
	}{
		{
			name: "check google_container_cluster with enable_legacy_abac set to true",
			source: `
resource "google_container_cluster" "gke" {
	enable_legacy_abac = "true"
	
}`,
			mustIncludeResultCode: checks.GkeAbacEnabled,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			results := scanSource(test.source)
			assertCheckCode(t, test.mustIncludeResultCode, test.mustExcludeResultCode, results)
		})
	}

}
