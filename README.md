# k8s auto ingress operator

为 srv 和 web 开头的 service 创建对应的 ingress

域名规则: `<serviceName>---<namespace>.<rootDomain>`


```bash
kgs

NAME                  TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
my-cs                 ClusterIP   10.43.118.186   <none>        5678/TCP   136m
srv-my-cs             ClusterIP   10.43.19.146    <none>        5678/TCP   136m
web-my-cs             ClusterIP   10.43.1.82      <none>        5678/TCP   134m


kg ing

NAME                  CLASS    HOSTS                                                   ADDRESS   PORTS     AGE
srv-my-cs--tangx-in   <none>   srv-my-cs---k8s-auto-ingress-operator-system.tangx.in             80, 443   88m
web-my-cs--tangx-in   <none>   web-my-cs---k8s-auto-ingress-operator-system.tangx.in             80, 443   9m17s
```


## 发布配置

1. 安装控制器

```
kubectl apply -f release/k8s-auto-ingress-operator.yml
```

2. 创建域名规则

```bash
kubectl apply -f deploy/tangx-in-ingresses.yml
```

配置文件如下

```yaml
# tangx-in-ingresses.yml
apiVersion: network.tangx.in/v1
kind: AutoIngress
metadata:
  name: tangx-in
  namespace: k8s-auto-ingress-operator-system

spec:
  rootDomain: tangx.in
  # tlsSecretName: "wild-tangx-in"
```


## 遗留问题

控制器启动时会获取所有的 service 。 如果这个时候没有 **域名规则** ， 将不会创建 ingress 规则。

1. 发布控制器
2. 发布规则
3. **删除控制器 pod， 重新ingess**


