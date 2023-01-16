package apps

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		},
	}
}

func appRedirectUrisResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).AppClient

	appId := d.Get("app_object_id").(string)

	app, _, err := client.Get(ctx, appId, odata.Query{})
	if err != nil {
		return tf.ErrorDiagF(err, "Application not found, app_id(%s)", appId)
	}

	if !d.HasChange("redirect_uris") {
		return appRedirectUrisResourceRead(ctx, d, meta)
	}

	toleranceOverride := d.Get("tolerance_override").(bool)

	oldRaw, newRaw := d.GetChange("redirect_uris")
	oldWeb, oldPub, oldSpa := classifyUrls(oldRaw)
	newWeb, newPub, newSpa := classifyUrls(newRaw)

	var webUrls = deRefSlice(app.Web.RedirectUris)
	webUrls, isOverrode := applyUrls(webUrls, oldWeb, newWeb)
	if !toleranceOverride && isOverrode {
		return diag.Errorf("Conflict detected: Some urls are already exist in target application, It may occur resource ownership conflict. If you want ignore this error, enable `tolerance_override` to true")
	}
	app.Web.RedirectUris = toRefSlice(webUrls)

	var pubUrls = deRefSlice(app.PublicClient.RedirectUris)
	pubUrls, isOverrode = applyUrls(pubUrls, oldPub, newPub)
	if !toleranceOverride && isOverrode {
		return diag.Errorf("Conflict detected: Some urls are already exist in target application, It may occur resource ownership conflict. If you want ignore this error, enable `tolerance_override` to true")
	}
	app.PublicClient.RedirectUris = toRefSlice(pubUrls)

	var spaUrls = deRefSlice(app.Spa.RedirectUris)
	spaUrls, isOverrode = applyUrls(spaUrls, oldSpa, newSpa)
	if !toleranceOverride && isOverrode {
		return diag.Errorf("Conflict detected: Some urls are already exist in target application, It may occur resource ownership conflict. If you want ignore this error, enable `tolerance_override` to true")
	}
	app.Spa.RedirectUris = toRefSlice(spaUrls)

	if _, err = client.Update(ctx, *app); err != nil {
		return tf.ErrorDiagF(err, "Application update failed, app_id(%s)", appId)
	}

	// Check Application updated to desired state
	updatedApp, _, err := client.Get(ctx, appId, odata.Query{})
	if err != nil {
		return tf.ErrorDiagF(err, "Application not found, app_id(%s)", appId)
	}
	if !equalsSlice(deRefSlice(app.Web.RedirectUris), deRefSlice(updatedApp.Web.RedirectUris)) {
		return diag.Errorf("Updated Application Web urls are not desired value, desired: %v, actual: %v", *app.Web.RedirectUris, *updatedApp.Web.RedirectUris)
	}
	if !equalsSlice(deRefSlice(app.PublicClient.RedirectUris), deRefSlice(updatedApp.PublicClient.RedirectUris)) {
		return diag.Errorf("Updated Application InstalledClient urls are not desired value, desired: %v, actual: %v", *app.PublicClient.RedirectUris, *updatedApp.PublicClient.RedirectUris)
	}
	if !equalsSlice(deRefSlice(app.Spa.RedirectUris), deRefSlice(updatedApp.Spa.RedirectUris)) {
		return diag.Errorf("Updated Application Spa urls are not desired value, desired: %v, actual: %v", *app.Spa.RedirectUris, *updatedApp.Spa.RedirectUris)
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
	client := meta.(*clients.Client).AppClient

	appId := d.Get("app_object_id").(string)

	app, _, err := client.Get(ctx, appId, odata.Query{})
	if err != nil {
		return tf.ErrorDiagF(err, "Cannot find target Application app_id(%s)", appId)
	}

	rawUris := d.Get("redirect_uris")
	oldWeb, oldPub, oldSpa := classifyUrls(rawUris)

	var webUrls = deRefSlice(app.Web.RedirectUris)
	webUrls, _ = applyUrls(webUrls, oldWeb, nil)
	app.Web.RedirectUris = toRefSlice(webUrls)

	var pubUrls = deRefSlice(app.PublicClient.RedirectUris)
	pubUrls, _ = applyUrls(pubUrls, oldPub, nil)
	app.PublicClient.RedirectUris = toRefSlice(pubUrls)

	var spaUrls = deRefSlice(app.Spa.RedirectUris)
	spaUrls, _ = applyUrls(spaUrls, oldSpa, nil)
	app.Spa.RedirectUris = toRefSlice(spaUrls)

	if _, err = client.Update(ctx, *app); err != nil {
		return tf.ErrorDiagF(err, "Application update failed, app_id(%s)", appId)
	}

	// Check Application updated to desired state
	updatedApp, _, err := client.Get(ctx, appId, odata.Query{})
	if err != nil {
		return tf.ErrorDiagF(err, "Application not found, app_id(%s)", appId)
	}
	if !equalsSlice(deRefSlice(app.Web.RedirectUris), deRefSlice(updatedApp.Web.RedirectUris)) {
		return diag.Errorf("Updated Application Web urls are not desired value, desired: %v, actual: %v", *app.Web.RedirectUris, *updatedApp.Web.RedirectUris)
	}
	if !equalsSlice(deRefSlice(app.PublicClient.RedirectUris), deRefSlice(updatedApp.PublicClient.RedirectUris)) {
		return diag.Errorf("Updated Application InstalledClient urls are not desired value, desired: %v, actual: %v", *app.PublicClient.RedirectUris, *updatedApp.PublicClient.RedirectUris)
	}
	if !equalsSlice(deRefSlice(app.Spa.RedirectUris), deRefSlice(updatedApp.Spa.RedirectUris)) {
		return diag.Errorf("Updated Application Spa urls are not desired value, desired: %v, actual: %v", *app.Spa.RedirectUris, *updatedApp.Spa.RedirectUris)
	}

	return nil
}

func classifyUrls(rawValue any) (web []string, pub []string, spa []string) {
	if rawValue == nil {
		return
	}

	for _, rawUrl := range rawValue.([]any) {
		urlSpec := rawUrl.(map[string]any)
		url := urlSpec["url"].(string)
		switch urlSpec["type"] {
		case "Web":
			web = append(web, url)
		case "InstalledClient":
			pub = append(pub, url)
		case "Spa":
			spa = append(spa, url)
		}
	}

	return
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

func deRefSlice(s *[]string) []string {
	if s == nil {
		return nil
	}
	return *s
}

func toRefSlice(s []string) *[]string {
	if s == nil {
		return nil
	}
	return &s
}
