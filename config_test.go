package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_DefaultAWSParameterStoreNodeIsSet(t *testing.T) {
	c := NewConfig()

	assert.Equal(t, "/dv/github-aws-credential-json", c.aws_parameter)
}

func TestConfig_DefaultGithubAPIUserIsSet(t *testing.T) {
	c := NewConfig()

	assert.Equal(t, "FKT-dv-jenkins", c.github_user)
}

func TestConfig_ValidateEmptyAWSProfileRaisesError(t *testing.T) {
	c := NewConfig()

	assert.Error(t, c.validate())
}

func TestConfig_ValidateNeedsProfileAndRepositoryToPass(t *testing.T) {
	c := NewConfig()
	c.github_repo = "some repo"
	c.aws_profile = "some profile"

	assert.NoError(t, c.validate())
}
