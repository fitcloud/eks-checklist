# NET-009 Endpoint 대신 EndpointSlices 사용

## Meaning
Kubernetes 1.16부터 도입된 EndpointSlices는 기존의 Endpoints 리소스를 개선한 것으로, 더 많은 Pod와 서비스를 효율적으로 관리할 수 있도록 돕습니다. EndpointSlices는 작고 효율적인 단위로 Pod 정보를 관리하여 성능을 향상시키고, 대규모 클러스터에서의 리소스 소모를 줄입니다.

## Impact
- 효율성 부족: Endpoint는 대규모 클러스터에서 성능 문제를 일으킬 수 있음
- 확장성 한계: 많은 Pod가 있을 경우 Endpoint 방식으로는 관리가 어려움

## Diagnosis
클러스터에서 EndpointSlices가 활성화되어 있는지 확인하고, Endpoint 대신 EndpointSlices를 사용하도록 설정합니다.

EndpointSlices 활성화 확인

```bash
kubectl get endpointslices -A
```

EndpointSlices가 사용되지 않으면, 해당 기능을 사용하도록 설정합니다. 
Kubernetes 1.16 이상에서는 기본적으로 활성화되어 있으므로, EndpointSlices가 사용되지 않으면, API 서버에서 설정을 확인합니다.

## Mitigation
대부분의 경우 Kubernetes 1.16 이상에서는 기본적으로 EndpointSlices를 사용하므로, 특별히 별도의 설정 없이 클러스터의 서비스가 EndpointSlices를 활용합니다.

[Kubenetes EndpointSlices](https://kubernetes.io/ko/docs/concepts/services-networking/endpoint-slices/)