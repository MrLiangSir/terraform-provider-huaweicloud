package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCdnDomainFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	hcCdnClient, err := cfg.HcCdnV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CDN v2 client: %s", err)
	}

	requestOpts := &model.ShowDomainDetailByNameRequest{
		DomainName:          state.Primary.Attributes["name"],
		EnterpriseProjectId: utils.StringIgnoreEmpty(state.Primary.Attributes["enterprise_project_id"]),
	}
	return hcCdnClient.ShowDomainDetailByName(requestOpts)
}

func TestAccCdnDomain_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckCDN(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.active", "1"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin", "100.254.53.75"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin_type", "ipaddr"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.http_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.https_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccCdnDomain_cache,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cache_settings.0.rules.0.rule_type", "all"),
					resource.TestCheckResourceAttr(resourceName, "cache_settings.0.rules.0.ttl", "180"),
					resource.TestCheckResourceAttr(resourceName, "cache_settings.0.rules.0.ttl_type", "d"),
					resource.TestCheckResourceAttr(resourceName, "cache_settings.0.rules.0.priority", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
				),
			},
			{
				Config: testAccCdnDomain_retrievalHost,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "sources.0.retrieval_host", "customize.test.huaweicloud.com"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.http_port", "8001"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.https_port", "8002"),
				),
			},
			{
				Config: testAccCdnDomain_standby,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "sources.0.active", "1"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin", "14.215.177.39"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin_type", "ipaddr"),
					resource.TestCheckResourceAttr(resourceName, "sources.1.active", "0"),
					resource.TestCheckResourceAttr(resourceName, "sources.1.origin", "220.181.28.52"),
					resource.TestCheckResourceAttr(resourceName, "sources.1.origin_type", "ipaddr"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCDNDomainImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id",
				},
			},
		},
	})
}

var testAccCdnDomain_basic = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

  configs {
    origin_protocol = "http"
  }

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
    http_port   = 80
    https_port  = 443
  }

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_cache = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

  configs {
    origin_protocol = "http"
  }

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
    http_port   = 80
    https_port  = 443
  }

  cache_settings {
    rules {
      rule_type = 0
      ttl       = 180
      ttl_type  = 4
      priority  = 2
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_retrievalHost = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

  configs {
    origin_protocol = "http"
  }

  sources {
    active         = 1
    origin         = "100.254.53.75"
    origin_type    = "ipaddr"
    retrieval_host = "customize.test.huaweicloud.com"
    http_port      = 8001
    https_port     = 8002
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_standby = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

  sources {
    active      = 1
    origin      = "14.215.177.39"
    origin_type = "ipaddr"
  }
  sources {
    active      = 0
    origin      = "220.181.28.52"
    origin_type = "ipaddr"
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

// Prepare the HTTPS certificate before running this test case
func TestAccCdnDomain_configHttpSettings(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckCDN(t)
			acceptance.TestAccPreCheckCERT(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_configHttpSettings,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ipv6_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.range_based_retrieval_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.certificate_name", "terraform-test"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.https_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.http2_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.cache_url_parameter_filter.0.type", "ignore_url_params"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.name", "test-name"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.status", "on"),
				),
			},
		},
	})
}

var testAccCdnDomain_configHttpSettings = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol               = "http"
    ipv6_enable                   = true
    range_based_retrieval_enabled = "true"

    https_settings {
      certificate_name = "terraform-test"
      certificate_body = file("%s")
      http2_enabled    = true
      https_enabled    = true
      private_key      = file("%s")
    }

    cache_url_parameter_filter {
      type = "ignore_url_params"
    }

    retrieval_request_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    http_response_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    url_signing {
      enabled = false
    }

    compress {
      enabled = false
    }

    force_redirect {
      enabled = true
      type    = "http"
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME, acceptance.HW_CDN_CERT_PATH, acceptance.HW_CDN_PRIVATE_KEY_PATH)

func TestAccCdnDomain_configs(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckCDN(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_configs,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ipv6_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.range_based_retrieval_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.cache_url_parameter_filter.0.type", "ignore_url_params"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.name", "test-name"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_frequency_limit.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_frequency_limit.0.qps", "1"),
				),
			},
			{
				Config: testAccCdnDomain_configsUpdate1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "follow"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ipv6_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.range_based_retrieval_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.cache_url_parameter_filter.0.type", "del_params"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.cache_url_parameter_filter.0.value", "test_value"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.name", "test-name-update"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.value", "test-val-update"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.action", "set"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_frequency_limit.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_frequency_limit.0.qps", "100000"),
				),
			},
			{
				Config: testAccCdnDomain_configsUpdate2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_frequency_limit.0.enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCDNDomainImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id",
				},
			},
		},
	})
}

var testAccCdnDomain_configs = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol               = "http"
    ipv6_enable                   = true
    range_based_retrieval_enabled = true

    cache_url_parameter_filter {
      type = "ignore_url_params"
    }

    retrieval_request_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    http_response_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    url_signing {
      enabled = false
    }

    compress {
      enabled = false
    }

    force_redirect {
      enabled = true
      type    = "http"
    }

    ip_frequency_limit {
      enabled = true
      qps     = 1
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_configsUpdate1 = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol               = "follow"
    ipv6_enable                   = false
    range_based_retrieval_enabled = false

    cache_url_parameter_filter {
      type  = "del_params"
      value = "test_value"
    }

    retrieval_request_header {
      name   = "test-name-update"
      value  = "test-val-update"
      action = "set"
    }

    http_response_header {
      name   = "Content-Disposition"
      value  = "test-val-update"
      action = "set"
    }

    url_signing {
      enabled = false
    }

    compress {
      enabled = false
    }

    force_redirect {
      enabled = true
      type    = "http"
    }

    ip_frequency_limit {
      enabled = true
      qps     = 100000
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_configsUpdate2 = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol               = "follow"
    ipv6_enable                   = false
    range_based_retrieval_enabled = false

    ip_frequency_limit {
      enabled = false
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

// This case is used to test fields that are only valid in the `wholeSite` scenario.
func TestAccCdnDomain_configTypeWholeSite(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckCDN(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_wholeSite,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "type", "wholeSite"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.timeout", "1"),
				),
			},
			{
				Config: testAccCdnDomain_wholeSiteUpdate1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "type", "wholeSite"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.timeout", "300"),
				),
			},
			{
				Config: testAccCdnDomain_wholeSiteUpdate2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "type", "wholeSite"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCDNDomainImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id",
				},
			},
		},
	})
}

var testAccCdnDomain_wholeSite = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "wholeSite"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol = "http"
    ipv6_enable     = true

    websocket {
      enabled = true
      timeout = 1
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_wholeSiteUpdate1 = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "wholeSite"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol = "http"
    ipv6_enable     = true

    websocket {
      enabled = true
      timeout = 300
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_wholeSiteUpdate2 = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "wholeSite"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol = "http"
    ipv6_enable     = true

    websocket {
      enabled = false
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

// testCDNDomainImportState use to return an ID using `name`
func testCDNDomainImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		return rs.Primary.Attributes["name"], nil
	}
}
