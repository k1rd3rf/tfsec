package tfsec

import (
	"testing"

	"github.com/k1rd3rf/tfsec/internal/app/tfsec/scanner"

	"github.com/k1rd3rf/tfsec/internal/app/tfsec/checks"
)

func Test_AWSUnencryptedS3Bucket(t *testing.T) {

	var tests = []struct {
		name                  string
		source                string
		mustIncludeResultCode scanner.RuleID
		mustExcludeResultCode scanner.RuleID
	}{
		{
			name: "check no server_side_encryption_configuration aws_s3_bucket",
			source: `
resource "aws_s3_bucket" "my-bucket" {
	
}`,
			mustIncludeResultCode: checks.AWSUnencryptedS3Bucket,
		},
		{
			name: "check no server_side_encryption_configuration aws_s3_bucket",
			source: `
resource "aws_s3_bucket" "my-bucket" {
  bucket = "mybucket"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "arn"
        sse_algorithm     = "aws:kms"
      }
    }
  }
}`,
			mustExcludeResultCode: checks.AWSUnencryptedS3Bucket,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			results := scanSource(test.source)
			assertCheckCode(t, test.mustIncludeResultCode, test.mustExcludeResultCode, results)
		})
	}

}
