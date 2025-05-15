# SEC-006 멀티 태넌시 적용 유무

## Meaning
멀티 태넌시는 하나의 클러스터에서 여러 사용자가 독립적으로 자원과 환경을 사용할 수 있도록 하는 방식입니다. 소프트 멀티 테넌시는 네임스페이스, RBAC, NetworkPolicy 등을 통해 소프트웨어적으로 자원을 격리하고, 하드 멀티 테넌시는 물리적인 격리를 강조하며, 클러스터를 분리하거나 각 테넌트마다 고립된 네트워크 및 IAM 정책을 적용합니다.

EKS에서는 네임스페이스 기반의 소프트 멀티 테넌시를 구현할 수 있지만, 하드 멀티 테넌시를 적용하려면 클러스터를 여러 개 배포하고 VPC, IAM, 네트워크 수준에서 완전한 격리가 필요합니다.

## Impact
- 소프트 멀티 테넌시: 같은 클러스터 내에서 네임스페이스를 기반으로 격리되므로, 일부 자원은 여전히 다른 테넌트와 공유됩니다.
- 보안 위험: 잘못된 RBAC 설정, 네트워크 정책 미적용 등으로 인해 테넌트 간 정보 유출 가능성
- 자원 충돌: 네임스페이스 간 리소스 경쟁 가능
- 하드 멀티 테넌시: 클러스터를 완전히 분리하여 리소스와 보안을 완전하게 격리
- 운영 복잡성 증가: 여러 클러스터를 관리해야 하므로, 관리 오버헤드와 비용이 증가

## Diagnosis
멀티 테넌시 적용 유무를 아래 항목으로 판단합니다

소프트 멀티 테넌시: 네임스페이스, RBAC, NetworkPolicy, IRSA, ResourceQuota 등을 사용하여 자원 격리 여부

하드 멀티 테넌시: 여러 EKS 클러스터 또는 VPC에 걸쳐 각 테넌트를 완전히 격리했는지 여부

클러스터 및 리소스 설정: 동일한 클러스터 내에서 자원 격리가 적절히 적용되었는지 확인

**example**

```bash
# 네임스페이스별 리소스 할당량 확인
kubectl get resourcequota -n <namespace>

# 네임스페이스 별 Role/RoleBinding 확인
kubectl get rolebindings -n <namespace>

# NetworkPolicy 확인
kubectl get networkpolicy -n <namespace>

# 여러 EKS 클러스터가 필요한 경우, 클러스터 간 리소스 격리 여부 점검
kubectl get nodes --context=<cluster-name>
```

## Mitigation
소프트 멀티 테넌시

네임스페이스를 사용하여 각 테넌트 격리

RBAC와 IAM Roles for Service Accounts(IRSA)을 통해 서비스 계정별 권한을 관리

NetworkPolicy로 네임스페이스 간 네트워크 격리

ResourceQuota와 LimitRange로 각 네임스페이스의 자원 사용 제한

하드 멀티 테넌시

각 테넌트마다 별도의 EKS 클러스터와 VPC를 생성하여 완전한 물리적 격리

IAM 정책을 클러스터별로 설정하여 권한 격리

**Before**

동일 클러스터 내에서 네임스페이스만으로 관리되는 환경

각 테넌트의 자원 격리가 불완전

**After**

소프트 멀티 테넌시 적용 시 각 테넌트 별 네임스페이스 분리 + RBAC + NetworkPolicy + ResourceQuota 적용

하드 멀티 테넌시 적용 시 각 테넌트 별 완전한 EKS 클러스터 및 VPC 격리

[IAM Roles for Service Accounts (IRSA) 공식 문서](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/iam-roles-for-service-accounts.html)

[EKS에서 하드 멀티 테넌시 구현하기 - AWS Well-Architected Framework](https://aws.amazon.com/ko/architecture/well-architected/?wa-lens-whitepapers.sort-by=item.additionalFields.sortDate&wa-lens-whitepapers.sort-order=desc&wa-guidance-whitepapers.sort-by=item.additionalFields.sortDate&wa-guidance-whitepapers.sort-order=desc)
