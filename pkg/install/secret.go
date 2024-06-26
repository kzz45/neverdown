package install

import (
	"encoding/base64"
	"encoding/json"

	"github.com/kzz45/neverdown/pkg/zaplogger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DockerSecret struct {
	Name     string
	Username string
	Password string
	Email    string
	Server   string
}

type DockerConfigJSON struct {
	Auths       DockerConfig      `json:"auths" datapolicy:"token"`
	HttpHeaders map[string]string `json:"HttpHeaders,omitempty" datapolicy:"token"`
}

type DockerConfig map[string]DockerConfigEntry

type DockerConfigEntry struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty" datapolicy:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth,omitempty" datapolicy:"token"`
}

func encodeDockerConfigFieldAuth(username, password string) string {
	fieldValue := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(fieldValue))
}

func handleDockerCfgJSONContent(username, password, email, server string) ([]byte, error) {
	dockerConfigAuth := DockerConfigEntry{
		Username: username,
		Password: password,
		Email:    email,
		Auth:     encodeDockerConfigFieldAuth(username, password),
	}
	dockerConfigJSON := DockerConfigJSON{
		Auths: map[string]DockerConfigEntry{server: dockerConfigAuth},
	}

	return json.Marshal(dockerConfigJSON)
}

func (b *installer) installDockerSecrets() {
	for _, dockerSecret := range b.opts.DockerSecrets {
		content, err := handleDockerCfgJSONContent(dockerSecret.Username, dockerSecret.Password, dockerSecret.Email, dockerSecret.Server)
		if err != nil {
			zaplogger.Sugar().Fatal(err)
		}
		for _, namespace := range b.supportedNamespaces() {
			sec := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      dockerSecret.Name,
				},
				Data: map[string][]byte{
					corev1.DockerConfigJsonKey: content,
				},
				Type: corev1.SecretTypeDockerConfigJson,
			}
			if _, err := b.clientSet.CoreV1().Secrets(namespace).Create(b.ctx, sec, metav1.CreateOptions{}); err != nil {
				if !errors.IsAlreadyExists(err) {
					zaplogger.Sugar().Fatal(err)
				}
				zaplogger.Sugar().Infof("Already exist Secrets:%s namespace:%s", dockerSecret.Name, namespace)
			} else {
				zaplogger.Sugar().Infof("Successful create Secrets:%s namespace:%s", dockerSecret.Name, namespace)
			}
		}
	}
}
