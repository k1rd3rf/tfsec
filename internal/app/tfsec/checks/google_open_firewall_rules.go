package checks

import (
	"fmt"
	"strings"

	"github.com/k1rd3rf/tfsec/internal/app/tfsec/scanner"

	"github.com/k1rd3rf/tfsec/internal/app/tfsec/parser"
)

// GoogleOpenInboundFirewallRule See https://github.com/liamg/tfsec#included-checks for check info
const GoogleOpenInboundFirewallRule scanner.RuleID = "GCP003"

// GoogleOpenOutboundFirewallRule See https://github.com/liamg/tfsec#included-checks for check info
const GoogleOpenOutboundFirewallRule scanner.RuleID = "GCP004"

func init() {
	scanner.RegisterCheck(scanner.Check{
		Code:           GoogleOpenInboundFirewallRule,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"google_compute_firewall"},
		CheckFunc: func(check *scanner.Check, block *parser.Block, _ *scanner.Context) []scanner.Result {

			if sourceRanges := block.GetAttribute("source_ranges"); sourceRanges != nil {

				if sourceRanges.Value().LengthInt() == 0 {
					return nil
				}

				for _, cidr := range sourceRanges.Value().AsValueSlice() {
					if strings.HasSuffix(cidr.AsString(), "/0") {
						return []scanner.Result{
							check.NewResult(
								fmt.Sprintf("Resource '%s' defines a fully open inbound firewall rule.", block.Name()),
								sourceRanges.Range(),
								scanner.SeverityWarning,
							),
						}
					}
				}

			}

			return nil
		},
	})

	scanner.RegisterCheck(scanner.Check{
		Code:           GoogleOpenOutboundFirewallRule,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"google_compute_firewall"},
		CheckFunc: func(check *scanner.Check, block *parser.Block, _ *scanner.Context) []scanner.Result {

			if destinationRanges := block.GetAttribute("destination_ranges"); destinationRanges != nil {

				if destinationRanges.Value().LengthInt() == 0 {
					return nil
				}

				for _, cidr := range destinationRanges.Value().AsValueSlice() {
					if strings.HasSuffix(cidr.AsString(), "/0") {
						return []scanner.Result{
							check.NewResult(
								fmt.Sprintf("Resource '%s' defines a fully open outbound firewall rule.", block.Name()),
								destinationRanges.Range(),
								scanner.SeverityWarning,
							),
						}
					}
				}

			}

			return nil
		},
	})
}
