# REL-017-DNS 캐시 적용

## Meaning
Codredns에 cache가 설정되어 잇는지 확인하여 DNS 쿼리 성능을 향상시키고 네트워크 리소스를 절약합니다.DNS cache은 DNS 응답을 일정 기간 동안 저장하여, 동일한 쿼리에 대해 반복적으로 DNS 서버에 요청하지 않고 캐시된 응답을 빠르게 반환할 수 있습니다.

## Impact
- DNS 응답 지연 증가: 캐시를 사용하지 않으면, DNS 요청이 있을 때마다 CoreDNS가 외부 DNS 서버나 다른 DNS 서버와 통신 해야 하기 때문에 네크워크 지연이 발생합니다.
- 쿼리 비용 증가: 외부 DNS 서버를 사용할 때, 특히 클라우드 환경에서 DNS 요청에 대한 비용이 발생하는 경우, 캐시를 사용하지 않으면 반복적인 DNS 쿼리가 많아져 쿼리 비용이 증가할 수 있습니다
- 성능 저하: 캐시를 사용하지 않으면 각 요청마다 DNS 서버에 쿼리를 보내야 하므로 서버 부하가 증가하고, 이는 시스템 전체의 성능 저하를 초래할 수 있습니다. 특히 요청이 많은 대규모 클러스터에서는 성능에 큰 영향을 미칠 수 있습니다

## Diagnosis
CoreDNS에 Cache 설정이 되어 있는지 확인하세요

```bash
kubectl get configmap coredns -n kube-system -o=jsonpath='{.data.Corefile}' | grep -i 'cache'
```

## Mitigation
CoreDNS에 Cache 설정을 해주세요
- CoreDNS 서버는 CoreDNS 구성 파일인 Core파일을 유지 관리하여 구성할 수 있습니다. 클러스터 관리자는 CoreDNS Corefile의 ConfigMap을 수정하여 해당 클러스터에 대한 DNS 서비스 검색 방식을 변경할 수 있습니다.

example
```bash
kubectl edit configmap coredns -n kube-system
```
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: coredns
  namespace: kube-system
data:
  Corefile: |
    .:53 {
        errors
        health {
            lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
            pods insecure
            fallthrough in-addr.arpa ip6.arpa
            ttl 30
        }
        prometheus :9153
        forward . /etc/resolv.conf
        cache 30
        loop
        reload
        loadbalance
    }    
```
[Amazon EKS 클러스터에서 DNS에 대한 CoreDNS 관리](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/managing-coredns.html)
[Customizing DNS Service](https://kubernetes.io/docs/tasks/administer-cluster/dns-custom-nameservers/)
[Corefile Explained](https://coredns.io/2017/07/23/corefile-explained/)