package helper

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewIngress(domain string, svc *corev1.Service) *netv1.Ingress {

	_domain := strings.ReplaceAll(domain, ".", "-")
	svcdomain := fmt.Sprintf("%s---%s.%s", svc.Name, svc.Namespace, domain)

	ing := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s--%s", svc.Name, _domain),
			Namespace: svc.Namespace,
			Labels:    svc.Labels,
			Annotations: map[string]string{
				"controller": "auto-ingress",
			},
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					Host: svcdomain,
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

	return ing
}

func ptrPathType(pt netv1.PathType) *netv1.PathType {
	return &pt
}
