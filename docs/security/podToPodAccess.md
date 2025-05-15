# SEC-009 Pod-to-Pod 접근 제어

## Meaning
기본적으로 Kubernetes 클러스터 내 모든 Pod들은 서로 자유롭게 통신할 수 있습니다. 이는 보안 위협 발생 시 lateral movement(수평 이동)를 가능하게 하며, 민감한 서비스가 외부에 노출되는 보안 취약점을 초래할 수 있습니다.

Pod-to-Pod 접근 제어는 이러한 무분별한 접근을 방지하고, 각 애플리케이션의 최소 권한 통신만 허용하는 네트워크 보안을 구성하는 것을 의미합니다.
이를 위해 Kubernetes의 NetworkPolicy를 사용해 특정 Pod 혹은 Namespace 간의 통신만 허용하는 정책을 정의합니다.

## Impact
- 무제한 통신 허용: 모든 Pod가 모든 Pod에 접근할 수 있어 보안 위협 발생 시 피해가 확산될 수 있음
- 서비스 간 데이터 노출: 인증/인가가 구현되지 않은 서비스가 노출될 수 있음
- 컴플라이언스 미준수: 보안 요구사항이 있는 환경에서 감사에 실패할 가능성 존재

## Diagnosis
아래 항목을 기준으로 Pod-to-Pod 접근 제어 적용 여부를 진단할 수 있습니다:

클러스터에서 NetworkPolicy가 활성화되어 있는지

각 네임스페이스 혹은 애플리케이션 단위로 NetworkPolicy가 적용되었는지

기본 거부(Deny all) 정책이 존재하는지 (policyTypes: Ingress, Egress)

사용하는 CNI(예: AWS VPC CNI, Calico 등)가 NetworkPolicy를 지원하는지 여부


```bash
# 클러스터에 존재하는 NetworkPolicy 목록 확인
kubectl get networkpolicy --all-namespaces

# 특정 네임스페이스의 정책 상세 확인
kubectl describe networkpolicy -n <namespace> <policy-name>
```

다음과 같은 기본 차단 정책이 적용되었는지 확인

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all-ingress
  namespace: example-namespace
spec:
  podSelector: {}
  policyTypes:
    - Ingress
```

## Mitigation
Pod 간 통신 제어를 위해 다음과 같은 작업을 수행해야 합니다

CNI 플러그인 확인

AWS VPC CNI는 기본적으로 NetworkPolicy를 지원하지 않음

Calico, Cilium 등을 설치하여 정책 기반 제어 가능

기본 차단 정책 적용

각 네임스페이스에 기본적으로 Ingress/Egress 모두 차단

허용 목록 작성

최소 권한 원칙에 따라 필요한 Pod 간의 통신만 명시적으로 허용

예: app=A가 app=B로 통신 가능하도록 설정

**Before**

모든 Pod가 서로 접근 가능

보안 위협 발생 시 내부 확산 가능

**After**

기본 거부 정책 + 최소 허용 정책 적용

민감한 서비스 보호 및 lateral movement 차단
