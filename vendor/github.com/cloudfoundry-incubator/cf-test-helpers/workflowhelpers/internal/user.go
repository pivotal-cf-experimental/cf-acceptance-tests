package internal

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/internal"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

type TestUser struct {
	username       string
	password       string
	cmdStarter     internal.Starter
	timeout        time.Duration
	shouldKeepUser bool
}

func NewTestUser(config config.Config, cmdStarter internal.Starter) *TestUser {
	node := ginkgoconfig.GinkgoConfig.ParallelNode
	timeTag := time.Now().Format("2006_01_02-15h04m05.999s")

	var regUser, regUserPass string
	regUser = fmt.Sprintf("%s-USER-%d-%s", config.NamePrefix, node, timeTag)
	regUserPass = "meow"

	if config.UseExistingUser {
		regUser = config.ExistingUser
		regUserPass = config.ExistingUserPassword
	}

	if config.ConfigurableTestPassword != "" {
		regUserPass = config.ConfigurableTestPassword
	}

	return &TestUser{
		username:       regUser,
		password:       regUserPass,
		cmdStarter:     cmdStarter,
		timeout:        config.ScaledTimeout(1 * time.Minute),
		shouldKeepUser: config.ShouldKeepUser,
	}
}

func NewAdminUser(config config.Config, cmdStarter internal.Starter) *TestUser {
	return &TestUser{
		username:   config.AdminUser,
		password:   config.AdminPassword,
		cmdStarter: cmdStarter,
	}
}

func (user *TestUser) Create() {
	session := internal.Cf(user.cmdStarter, "create-user", user.username, user.password)
	EventuallyWithOffset(1, session, user.timeout).Should(Exit())
	if session.ExitCode() != 0 {
		ExpectWithOffset(1, session.Out).Should(Say("scim_resource_already_exists"))
	}
}

func (user *TestUser) Destroy() {
	session := internal.Cf(user.cmdStarter, "delete-user", "-f", user.username)
	EventuallyWithOffset(1, session, user.timeout).Should(Exit(0))
}

func (user *TestUser) Username() string {
	return user.username
}

func (user *TestUser) Password() string {
	return user.password
}

func (user *TestUser) ShouldRemain() bool {
	return user.shouldKeepUser
}
