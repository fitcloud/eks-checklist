# NET-008 kube-proxy에 IPVS 모드 적용

## Meaning
IP 가상 서버(IPVS) 모드의 EKS는 기존 iptables 모드에서 실행되는 kube-proxy를 사용하여 대규모 클러스터에서 1,000개 이상의 서비스를 제공할 때 발생하는 네트워크 지연 문제를 해결합니다. 이 성능 문제는 각 패킷에 대해 iptables 패킷 필터링 규칙을 순차적으로 처리하는 방식에서 비롯됩니다. 이를 해결하려면 클러스터를 IPVS 모드에서 kube-proxy를 실행하도록 구성할 수 있습니다

## Impact
- 성능 저하: iptables는 패킷 필터링 규칙을 순차적으로 처리하는 방식이기 때문에, 서비스 수가 많아질수록 처리해야 할 규칙이 증가합니다.
- 비효율적인 로드 밸런싱: iptables는 기본적으로 패킷을 순차적으로 검사하기 때문에 로드 밸런싱 성능이 떨어질 수 있습니다
- 패킷 처리 속도 저하: 패킷 필터링 규칙을 순차적으로 처리하는 방식 때문에, 패킷 처리 속도가 느려져 네트워크 지연(latency)이 발생할 수 있습니다.
- 확장성: iptalbes 모드는 서비스 규모가 커질수록 성능 저하가 발생하는 반면면 IPVS 모드는 커널 레벨에서 로드밸런싱을 수행해 높은 성능과 확장성을 제공하므로, 대규모 클러스터 환경에서 더 효과적으로 동작합니다

## Diagnosis
kube-proxy-config mode를 확인하세요
```bash
kubectl get cm -n kube-system kube-proxy-config -o json | jq -r '.data.config' | grep -E 'mode'
```

## Mitigation
kube-proxy-config mode를 IPVS로 변경해주세요.

**example**

EKS Add-on kube-proxy 업데이트
```bash
aws eks update-addon --cluster-name $CLUSTER_NAME --addon-name kube-proxy \
  --configuration-values '{"ipvs": {"scheduler": "rr"}, "mode": "ipvs"}' \
  --resolve-conflicts OVERWRITE
```
EKS 클러스터 kube-proxy-config ConfigMap 수정
```bash
kubectl -n kube-system edit cm kube-proxy-config
```
```yaml
  iptables:
    masqueradeAll: false
    masqueradeBit: 14
    minSyncPeriod: 0s
    syncPeriod: 30s
  ipvs:
    excludeCIDRs: null
    minSyncPeriod: 0s
    scheduler: "rr"
    syncPeriod: 30s
  kind: KubeProxyConfiguration
  metricsBindAddress: 0.0.0.0:10249
  mode: "ipvs"
  nodePortAddresses: null
  oomScoreAdj: -998
  portRange: ""
  udpIdleTimeout: 250ms
```

[Kubenetes Kube-proxy란](https://kubernetes.io/ko/docs/reference/command-line-tools-reference/kube-proxy/)
[IPVS 모드에서 kube-proxy 실행](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/ipvs.html)