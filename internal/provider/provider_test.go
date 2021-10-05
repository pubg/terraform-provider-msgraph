package provider

import (
	"context"
	"fmt"
	"github.com/manicminer/hamilton/auth"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"terraform-provider-msgraph/internal/clients"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func TestAccProvider_cliAuth(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		return
	}

	provider := Provider()
	ctx := context.Background()

	// Support only Azure CLI authentication
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		environment, aadEnvironment := environment(d.Get("environment").(string))

		aadBuilder := &authentication.Builder{
			Environment: aadEnvironment,
			TenantID:    d.Get("tenant_id").(string),
			TenantOnly:  true,

			SupportsAzureCliToken: true,
		}

		authConfig := &auth.Config{
			Environment: environment,
			TenantID:    d.Get("tenant_id").(string),

			EnableAzureCliToken: true,
		}

		return buildClient(ctx, provider, authConfig, aadBuilder, "", true)
	}

	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d != nil && d.HasError() {
		t.Fatalf("err: %+v", d)
	}

	if errs := testCheckProvider(provider); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestAccProvider_clientCertificateAuth(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		return
	}

	provider := Provider()
	ctx := context.Background()

	// Support only Service Principal authentication (certificate or secret)
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		environment, aadEnvironment := environment(d.Get("environment").(string))

		aadBuilder := &authentication.Builder{
			Environment: aadEnvironment,
			TenantID:    d.Get("tenant_id").(string),
			ClientID:    d.Get("client_id").(string),
			TenantOnly:  true,

			SupportsClientCertAuth: true,
			ClientCertPassword:     d.Get("client_certificate_password").(string),
			ClientCertPath:         d.Get("client_certificate_path").(string),
		}

		authConfig := &auth.Config{
			Environment: environment,
			TenantID:    d.Get("tenant_id").(string),
			ClientID:    d.Get("client_id").(string),

			EnableClientCertAuth: true,
			ClientCertPath:       d.Get("client_certificate_path").(string),
			ClientCertPassword:   d.Get("client_certificate_password").(string),
		}

		return buildClient(ctx, provider, authConfig, aadBuilder, "", true)
	}

	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d != nil && d.HasError() {
		t.Fatalf("err: %+v", d)
	}

	if errs := testCheckProvider(provider); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestAccProvider_clientSecretAuth(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		return
	}

	provider := Provider()
	ctx := context.Background()

	// Support only Service Principal authentication (certificate or secret)
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		environment, aadEnvironment := environment(d.Get("environment").(string))

		aadBuilder := &authentication.Builder{
			Environment: aadEnvironment,
			TenantID:    d.Get("tenant_id").(string),
			ClientID:    d.Get("client_id").(string),
			TenantOnly:  true,

			SupportsClientSecretAuth: true,
			ClientSecret:             d.Get("client_secret").(string),
		}

		authConfig := &auth.Config{
			Environment: environment,
			TenantID:    d.Get("tenant_id").(string),
			ClientID:    d.Get("client_id").(string),

			EnableClientSecretAuth: true,
			ClientSecret:           d.Get("client_secret").(string),
		}

		return buildClient(ctx, provider, authConfig, aadBuilder, "", true)
	}

	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d != nil && d.HasError() {
		t.Fatalf("err: %+v", d)
	}

	if errs := testCheckProvider(provider); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func testCheckProvider(provider *schema.Provider) (errs []error) {
	client := provider.Meta().(*clients.Client)

	if client.Environment.AzureADEndpoint == "" {
		errs = append(errs, fmt.Errorf("AzureADEndpoint was empty in client.Environment"))
	}

	if client.Environment.MsGraph.Endpoint == "" {
		errs = append(errs, fmt.Errorf("MsGraphEndpoint was empty in client.Environment"))
	}

	if client.TenantID == "" {
		errs = append(errs, fmt.Errorf("client.TenantID was empty"))
	}

	if client.ClientID == "" {
		errs = append(errs, fmt.Errorf("client.ClientID was empty"))
	}

	if client.Claims.TenantId == "" {
		errs = append(errs, fmt.Errorf("TenantId was not populated in client.Claims"))
	}

	if client.Claims.ObjectId == "" {
		errs = append(errs, fmt.Errorf("ObjectId was not populated in client.Claims"))
	}

	return
}