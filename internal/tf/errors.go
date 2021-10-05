package tf

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ErrorDiagF(err error, format string, a ...interface{}) diag.Diagnostics {
	return ErrorDiagPathF(err, "", format, a...)
}

func ErrorDiagPathF(err error, attr string, summary string, a ...interface{}) diag.Diagnostics {
	d := diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf(summary, a...),
	}
	if err != nil {
		d.Detail = err.Error()
	}
	if attr != "" {
		d.AttributePath = cty.Path{cty.GetAttrStep{Name: attr}}
	}
	return diag.Diagnostics{d}
}
