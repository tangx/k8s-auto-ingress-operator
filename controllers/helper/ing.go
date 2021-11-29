package helper

import (
	"fmt"
	"strings"

	v1 "github.com/tangx/k8s-auto-ingress-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewIngress(op v1.AutoIngress, svc *corev1.Service) *netv1.Ingress {
	domain := op.Spec.RootDomain

	_domain := strings.ReplaceAll(domain, ".", "-")

	host := fmt.Sprintf("%s---%s.%s", svc.Name, svc.Namespace, domain)
	ingname := fmt.Sprintf("%s--%s", svc.Name, _domain)

	ing := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ingname,
			Namespace:   svc.Namespace,
			Labels:      svc.Labels,
			Annotations: annotations(&op),
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: ptrPathType(netv1.PathTypePrefix),
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: svc.Name,
											Port: netv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if op.Spec.TlsSecretName != nil {
		ing.Spec.TLS = []netv1.IngressTLS{
			{
				Hosts: []string{
					host,
				},
				SecretName: *op.Spec.TlsSecretName,
			},
		}
	}

	return ing
}

func ptrPathType(pt netv1.PathType) *netv1.PathType {
	return &pt
}

func annotations(objs ...client.Object) map[string]string {

	annos := make(map[string]string)
	for _, obj := range objs {

		for k, v := range obj.GetAnnotations() {
			annos[k] = v
		}
	}

	return annos

}
