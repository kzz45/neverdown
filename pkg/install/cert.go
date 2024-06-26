package install

import (
	"context"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CertificateOption struct {
	CaCrt     []byte
	ServerCrt []byte
	ServerKey []byte
}

const (
	KubeApiCertificateName = "openx-root-certs"
	KubeApiDashboardName   = "openx-dashboard-config-json"
)

func (b *installer) installCertificateConfigMap() {
	for _, namespace := range b.supportedNamespaces() {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      KubeApiCertificateName,
			},
			Immutable: nil,
			Data: map[string]string{
				"ca.crt":     string(b.opts.CertificateOption.CaCrt),
				"server.crt": string(b.opts.CertificateOption.ServerCrt),
				"server.key": string(b.opts.CertificateOption.ServerKey),
			},
			BinaryData: nil,
		}
		if _, err := b.clientSet.CoreV1().ConfigMaps(namespace).Create(context.Background(), cm, metav1.CreateOptions{}); err != nil {
			if !errors.IsAlreadyExists(err) {
				zaplogger.Sugar().Fatal(err)
			}
			zaplogger.Sugar().Infof("Already exist ConfigMaps:%s namespace:%s", cm.Name, cm.Namespace)
		} else {
			zaplogger.Sugar().Infof("Successful create ConfigMaps:%s namespace:%s", cm.Name, cm.Namespace)
		}
	}
}
