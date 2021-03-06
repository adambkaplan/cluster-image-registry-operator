package resource

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	coreset "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/openshift/cluster-image-registry-operator/pkg/defaults"
)

var _ Mutator = &generatorPullSecret{}

type generatorPullSecret struct {
	client    coreset.CoreV1Interface
	namespace string
}

func newGeneratorPullSecret(client coreset.CoreV1Interface) *generatorPullSecret {
	return &generatorPullSecret{
		client:    client,
		namespace: defaults.ImageRegistryOperatorNamespace,
	}
}

func (gs *generatorPullSecret) Type() runtime.Object {
	return &corev1.Secret{}
}

func (gs *generatorPullSecret) GetGroup() string {
	return corev1.GroupName
}

func (gs *generatorPullSecret) GetResource() string {
	return "secrets"
}

func (gs *generatorPullSecret) GetNamespace() string {
	return gs.namespace
}

func (gs *generatorPullSecret) GetName() string {
	return defaults.InstallationPullSecret
}

func (gs *generatorPullSecret) expected() (runtime.Object, error) {
	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gs.GetName(),
			Namespace: gs.GetNamespace(),
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{},
	}

	orig, err := gs.client.Secrets("openshift-config").Get(
		context.TODO(), "pull-secret", metav1.GetOptions{},
	)
	if errors.IsNotFound(err) {
		return sec, nil
	} else if err != nil {
		return nil, err
	}

	sec.Data = orig.Data
	return sec, nil
}

func (gs *generatorPullSecret) Get() (runtime.Object, error) {
	return gs.client.Secrets(gs.GetNamespace()).Get(
		context.TODO(), gs.GetName(), metav1.GetOptions{},
	)
}

func (gs *generatorPullSecret) Create() (runtime.Object, error) {
	return commonCreate(gs, func(obj runtime.Object) (runtime.Object, error) {
		return gs.client.Secrets(gs.GetNamespace()).Create(
			context.TODO(), obj.(*corev1.Secret), metav1.CreateOptions{},
		)
	})
}

func (gs *generatorPullSecret) Update(o runtime.Object) (runtime.Object, bool, error) {
	return commonUpdate(gs, o, func(obj runtime.Object) (runtime.Object, error) {
		return gs.client.Secrets(gs.GetNamespace()).Update(
			context.TODO(), obj.(*corev1.Secret), metav1.UpdateOptions{},
		)
	})
}

func (gs *generatorPullSecret) Delete(opts metav1.DeleteOptions) error {
	return gs.client.Secrets(gs.GetNamespace()).Delete(
		context.TODO(), gs.GetName(), opts,
	)
}

func (g *generatorPullSecret) Owned() bool {
	return true
}
