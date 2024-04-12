package env

import "os"

const (
	DomainSource       = "DOMAIN_SOURCE"
	DomainTarget       = "DOMAIN_TARGET"
	DomainSourceRegion = "DOMAIN_SOURCE_REGION"
	DomainTargetRegion = "DOMAIN_TARGET_REGION"
)

func DomainRefactorRules() []string {
	s1, ok := os.LookupEnv(DomainSource)
	if !ok {
		return []string{}
	}
	s2, ok := os.LookupEnv(DomainTarget)
	if !ok {
		return []string{}
	}
	s3, ok := os.LookupEnv(DomainSourceRegion)
	if !ok {
		return []string{s1, s2}
	}
	s4, ok := os.LookupEnv(DomainTargetRegion)
	if !ok {
		return []string{s1, s2}
	}
	return []string{s1, s2, s3, s4}
}
