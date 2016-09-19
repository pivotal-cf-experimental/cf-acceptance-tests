package security_groups_test

import (
	"testing"
	"time"

	"github.com/cloudfoundry/cf-acceptance-tests/helpers/skip_messages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cf_config "github.com/cloudfoundry-incubator/cf-test-helpers/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
)

var (
	DEFAULT_TIMEOUT   = 30 * time.Second
	CF_PUSH_TIMEOUT   = 2 * time.Minute
	LONG_CURL_TIMEOUT = 2 * time.Minute
)

var (
	testSetup            *workflowhelpers.ReproducibleTestSuiteSetup
	config               cf_config.Config
	DEFAULT_MEMORY_LIMIT = "256M"
)

func TestApplications(t *testing.T) {
	RegisterFailHandler(Fail)

	config = cf_config.LoadConfig()

	if config.DefaultTimeout > 0 {
		DEFAULT_TIMEOUT = config.DefaultTimeout * time.Second
	}

	if config.CfPushTimeout > 0 {
		CF_PUSH_TIMEOUT = config.CfPushTimeout * time.Second
	}

	if config.LongCurlTimeout > 0 {
		LONG_CURL_TIMEOUT = config.LongCurlTimeout * time.Second
	}

	testSetup = workflowhelpers.NewTestSuiteSetup(config)

	BeforeSuite(func() {
		testSetup.Setup()
	})

	AfterSuite(func() {
		testSetup.Teardown()
	})

	componentName := "SecurityGroups"

	rs := []Reporter{}

	if config.ArtifactsDirectory != "" {
		helpers.EnableCFTrace(config, componentName)
		rs = append(rs, helpers.NewJUnitReporter(config, componentName))
	}

	RunSpecsWithDefaultAndCustomReporters(t, componentName, rs)
}

func SecurityGroupsDescribe(description string, callback func()) bool {
	BeforeEach(func() {
		if !config.IncludeSecurityGroups {
			Skip(skip_messages.SkipSecurityGroupsMessage)
		}
	})
	return Describe("[security_groups] "+description, callback)
}
