package repositories

import (
	`testing`

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRepositories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "postgres suite")
}
