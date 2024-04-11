package config

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

// LeaderElectionConfiguration defines the configuration of leader election
// clients for components that can run with leader election enabled.
type LeaderElectionConfiguration struct {
	// leaderElect enables a leader election client to gain leadership
	// before executing the main loop. Enable this when running replicated
	// components for high availability.
	LeaderElect bool

	// HealthzAdapterDuration
	HealthzAdapterDuration metav1.Duration

	// leaseDuration is the duration that non-leader candidates will wait
	// after observing a leadership renewal until attempting to acquire
	// leadership of a led but unrenewed leader slot. This is effectively the
	// maximum duration that a leader can be stopped before it is replaced
	// by another candidate. This is only applicable if leader election is
	// enabled.
	LeaseDuration metav1.Duration
	// renewDeadline is the interval between attempts by the acting master to
	// renew a leadership slot before it stops leading. This must be less
	// than or equal to the lease duration. This is only applicable if leader
	// election is enabled.
	RenewDeadline metav1.Duration
	// retryPeriod is the duration the clients should wait between attempting
	// acquisition and renewal of a leadership. This is only applicable if
	// leader election is enabled.
	RetryPeriod metav1.Duration
	// resourceLock indicates the resource object type that will be used to lock
	// during leader election cycles.
	ResourceLock string
	// resourceName indicates the name of resource object that will be used to lock
	// during leader election cycles.
	ResourceName string
	// resourceName indicates the namespace of resource object that will be used to lock
	// during leader election cycles.
	ResourceNamespace string

	// ReleaseOnCancel should be set true if the lock should be released
	// when the run context is cancelled. If you set this to true, you must
	// ensure all code guarded by this lease has successfully completed
	// prior to cancelling the context, or you may have two processes
	// simultaneously acting on the critical path.
	ReleaseOnCancel bool
}

func DefaultLeaderConfig() *LeaderElectionConfiguration {
	return &LeaderElectionConfiguration{
		LeaderElect:            true,
		HealthzAdapterDuration: metav1.Duration{Duration: 20 * time.Second},
		LeaseDuration:          metav1.Duration{Duration: 10 * time.Second},
		RenewDeadline:          metav1.Duration{Duration: 5 * time.Second},
		RetryPeriod:            metav1.Duration{Duration: 2 * time.Second},
		ResourceLock:           resourcelock.LeasesResourceLock,
		ResourceName:           "openx-controller",
		ResourceNamespace:      "kube-neverdown",
		ReleaseOnCancel:        true,
	}
}
