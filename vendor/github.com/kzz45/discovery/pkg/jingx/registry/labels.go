package registry

import (
	"fmt"
	"strings"

	"github.com/kzz45/discovery/pkg/zaplogger"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/validation"
)

const (
	JingxProject    = "jingx-project"
	JingxRepository = "jingx-repository"
)

// GetLabelSelector returns the LabelSelector of the metav1.ListOptions
func GetLabelSelector(in map[string]string) string {
	ls := labels.NewSelector()
	for k, v := range in {
		req, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			zaplogger.Sugar().Fatal(err)
		}
		ls = ls.Add(*req)
	}
	return ls.String()
}

func ValidateName(name string) error {
	if errs := validation.IsValidLabelValue(name); len(errs) != 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}
	return nil
}
