# REL-017 DNS 캐시 적용

## Meaning
CoreDNS에 캐시 설정 여부를 확인하여 DNS 쿼리 성능을 향상시키고, 네트워크 리소스를 절약합니다. DNS 캐시는 DNS 응답을 일정 기간 동안 저장하여, 동일한 쿼리에 대해 반복적으로 DNS 서버에 요청하지 않고 캐시된 응답을 빠르게 반환할 수 있습니다.

## Impact
- DNS 응답 지연 증가: 캐시 미설정 시, DNS 요청마다 외부 DNS 서버와 통신해야 하므로 네트워크 지연이 발생합니다.
- 쿼리 비용 증가: 외부 DNS 서버를 사용할 경우 반복적인 DNS 쿼리로 인해 비용이 증가할 수 있습니다.
- 성능 저하: 캐시 미설정 시, 각 요청마다 DNS 서버에 쿼리를 보내야 하므로 서버 부하가 증가하고 성능 저하를 초래할 수 있습니다.

## Diagnosis
다음 명령어를 통해 CoreDNS에 캐시 설정이 되어 있는지 확인합니다

```bash
kubectl get configmap coredns -n kube-system -o=jsonpath='{.data.Corefile}' | grep -i 'cache'
```

## Mitigation
CoreDNS에 캐시 설정을 추가하려면 Corefile을 수정해야 합니다. 다음과 같이 배포 시 config의 cache 설정을 추가합니다.

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