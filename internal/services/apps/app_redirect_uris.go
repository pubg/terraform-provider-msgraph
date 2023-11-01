package apps

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
	"github.com/pubg/terraform-provider-msgraph/internal/clients"
	"github.com/pubg/terraform-provider-msgraph/internal/tf"
)

func appRedirectUris() *schema.Resource {
	return &schema.Resource{
		CreateContext: appRedirectUrisResourceUpdate,
		UpdateContext: appRedirectUrisResourceUpdate,
		ReadContext:   appRedirectUrisResourceRead,
		DeleteContext: appRedirectUrisResourceDelete,

		Schema: map[string]*schema.Schema{
			"app_object_id": {
				Description:      "object_id of target Application",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},

			"redirect_uris": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Description:      "A reply URL, You can find restrictions and limitations are this document. https://learn.microsoft.com/en-us/azure/active-directory/develop/reply-url",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
						"type": {
							Description:      "Type of Redirect Url, One of [Web, InstalledClient, Spa]",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Web", "InstalledClient", "Spa"}, false)),
						},
					},
				},
			},

			"tolerance_override": {
				Description: "If some urls are already exist in target application, It may occur resource ownership conflict. If you want ignore this error, enable `tolerance_override` to true",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},

			"retry_count": {
				Description: "Retry count for update application",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
			},
		},
	}
}

type RedirectUrlType string

const (
	Web             = RedirectUrlType("Web")
	InstalledClient = RedirectUrlType("InstalledClient")
	Spa             = RedirectUrlType("Spa")
)

const UrlOverrodeDetectedMessage = "Conflict detected: Some urls are already exist in target application, It may occur resource ownership conflict. If you want ignore this error, enable `tolerance_override` to true"

func appRedirectUrisResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerClient := meta.(*clients.Client)
	client := providerClient.AppClient

	if providerClient.EnableResourceMutex {
		providerClient.ResourceMutex.Lock()
		defer providerClient.ResourceMutex.Unlock()
	}

	if !d.HasChange("redirect_uris") {
		return appRedirectUrisResourceRead(ctx, d, meta)
	}
	oldRaw, newRaw := d.GetChange("redirect_uris")
	appId := d.Get("app_object_id").(string)
	toleranceOverride := d.Get("tolerance_override").(bool)
	retryCount := d.Get("retry_count").(int)

	var diagnostics diag.Diagnostics
	for i := 0; i < retryCount; i++ {
		diagnostics = updateTerraformResource(ctx, appId, client, oldRaw, newRaw, toleranceOverride)
		if diagnostics == nil {
			break
		}
	}
	if diagnostics != nil {
		return diagnostics
	}

	return appRedirectUrisResourceRead(ctx, d, meta)
}

func appRedirectUrisResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).AppClient

	appId := d.Get("app_object_id").(string)

	_, _, err := client.Get(ctx, appId, odata.Query{})
	if err != nil {
		return tf.ErrorDiagF(err, "Cannot find target Application app_id(%s)", appId)
	}

	d.SetId(appId)
	return nil
}

func appRedirectUrisResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerClient := meta.(*clients.Client)
	client := providerClient.AppClient

	if providerClient.EnableResourceMutex {
		providerClient.ResourceMutex.Lock()
		defer providerClient.ResourceMutex.Unlock()
	}

	appId := d.Get("app_object_id").(string)
	rawUris := d.Get("redirect_uris")
	toleranceOverride := d.Get("tolerance_override").(bool)
	retryCount := d.Get("retry_count").(int)

	var diagnostics diag.Diagnostics
	for i := 0; i < retryCount; i++ {
		diagnostics = updateTerraformResource(ctx, appId, client, rawUris, nil, toleranceOverride)
		if diagnostics == nil {
			break
		}
	}
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func updateTerraformResource(ctx context.Context, appId string, client *msgraph.ApplicationsClient, oldRaw any, newRaw any, toleranceOverride bool) diag.Diagnostics {
	app, _, err := client.Get(ctx, appId, odata.Query{})
	if err != nil {
		return tf.ErrorDiagF(err, "Application not found, app_id(%s)", appId)
	}

	if isOverrode := applyUrlsToApplication(Web, app, oldRaw, newRaw); !toleranceOverride && isOverrode {
		return diag.Errorf(UrlOverrodeDetectedMessage)
	}
	if isOverrode := applyUrlsToApplication(InstalledClient, app, oldRaw, newRaw); !toleranceOverride && isOverrode {
		return diag.Errorf(UrlOverrodeDetectedMessage)
	}
	if isOverrode := applyUrlsToApplication(Spa, app, oldRaw, newRaw); !toleranceOverride && isOverrode {
		return diag.Errorf(UrlOverrodeDetectedMessage)
	}

	if _, err = client.Update(ctx, *app); err != nil {
		return tf.ErrorDiagF(err, "Application update failed, app_id(%s)", appId)
	}

	// Check Application updated to desired state
	updatedApp, _, err := client.Get(ctx, appId, odata.Query{})
	if err != nil {
		return tf.ErrorDiagF(err, "Application not found, app_id(%s)", appId)
	}
	if diagnostics := checkApplicationIsDesiredStateByType(Web, app, updatedApp); diagnostics != nil {
		return diagnostics
	}
	if diagnostics := checkApplicationIsDesiredStateByType(InstalledClient, app, updatedApp); diagnostics != nil {
		return diagnostics
	}
	if diagnostics := checkApplicationIsDesiredStateByType(Spa, app, updatedApp); diagnostics != nil {
		return diagnostics
	}
	return nil
}

func checkApplicationIsDesiredStateByType(urlType RedirectUrlType, desire *msgraph.Application, actual *msgraph.Application) diag.Diagnostics {
	desireUrls := getUrlsByTypeFromApplication(urlType, desire)
	actualUrls := getUrlsByTypeFromApplication(urlType, actual)
	if !equalsSlice(desireUrls, actualUrls) {
		return diag.Errorf("Updated Application %s urls are not desired value, desired: %v, actual: %v", urlType, desireUrls, actualUrls)
	}
	return nil
}

// Returns url overrode detected
func applyUrlsToApplication(urlType RedirectUrlType, app *msgraph.Application, oldRedirectUris any, newRedirectUris any) bool {
	appUrls := getUrlsByTypeFromApplication(urlType, app)
	oldUrls := getUrlsByTypeFromTerraform(urlType, oldRedirectUris)
	newUrls := getUrlsByTypeFromTerraform(urlType, newRedirectUris)

	urls, isOverrode := applyUrls(appUrls, oldUrls, newUrls)

	setUrlsByTypeToApplication(urlType, app, urls)
	return isOverrode
}

func getUrlsByTypeFromApplication(urlType RedirectUrlType, app *msgraph.Application) []string {
	if urlType == Web {
		if app.Web.RedirectUris == nil {
			return nil
		}
		return *app.Web.RedirectUris
	} else if urlType == InstalledClient {

		if app.PublicClient.RedirectUris == nil {
			return nil
		}
		return *app.PublicClient.RedirectUris
	} else {
		if app.Spa.RedirectUris == nil {
			return nil
		}
		return *app.Spa.RedirectUris
	}
}

func setUrlsByTypeToApplication(urlType RedirectUrlType, app *msgraph.Application, urls []string) {
	if urlType == Web {
		if urls == nil {
			app.Web.RedirectUris = nil
		} else {
			app.Web.RedirectUris = &urls
		}
	} else if urlType == InstalledClient {
		if urls == nil {
			app.PublicClient.RedirectUris = nil
		} else {
			app.PublicClient.RedirectUris = &urls
		}
	} else {
		if urls == nil {
			app.Spa.RedirectUris = nil
		} else {
			app.Spa.RedirectUris = &urls
		}
	}
}

func getUrlsByTypeFromTerraform(urlType RedirectUrlType, rawValue any) []string {
	if rawValue == nil {
		return nil
	}

	var result []string
	for _, rawUrl := range rawValue.([]any) {
		urlSpec := rawUrl.(map[string]any)
		if urlSpec["type"].(string) == string(urlType) {
			result = append(result, urlSpec["url"].(string))
		}
	}
	return result
}

func applyUrls(currentUrls []string, oldUrls []string, newUrls []string) (urls []string, overrode bool) {
	urlMap := map[string]any{}
	for _, url := range currentUrls {
		urlMap[url] = nil
	}
	for _, url := range oldUrls {
		delete(urlMap, url)
	}
	for _, url := range newUrls {
		if _, found := urlMap[url]; found {
			// Find overrode
			overrode = true
		}
		urlMap[url] = nil
	}

	for url, _ := range urlMap {
		urls = append(urls, url)
	}
	return
}

func equalsSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := map[string]any{}
	for _, v := range a {
		aMap[v] = nil
	}
	for _, v := range b {
		if _, found := aMap[v]; found {
			delete(aMap, v)
		} else {
			return false
		}
	}

	if len(aMap) != 0 {
		return false
	}

	return true
}
