# SEC-003 IRSA 또는 EKS Pod Identity 기반 권한 부여

## Meaning
EKS에서 워크로드(Pod)가 AWS 리소스에 접근할 수 있도록 권한을 부여할 때, 전통적으로는 노드 IAM Role을 사용했지만 이는 과도한 권한 부여의 문제가 있습니다.

IRSA(IAM Roles for Service Accounts) 또는 EKS Pod Identity는 각 서비스 계정 단위로 IAM Role을 부여할 수 있어 최소 권한 원칙을 따르는 보안적인 접근 방식입니다.

- 두 가지 방식
  - IRSA (IAM Roles for Service Accounts): ServiceAccount에 IAM Role을 연결해 AWS 리소스에 안전하게 접근
  - EKS Pod Identity: EKS에서 보다 간단하게 IAM Role을 Pod에 매핑할 수 있는 방식 (EKS 전용 기능, 향후 IRSA 대체 가능)

## Impact
IRSA나 EKS Pod Identity를 사용하지 않는 경우 다음과 같은 문제가 발생할 수 있습니다

- 노드 Role을 공유하여 불필요하게 많은 권한이 Pod에 전달됨
- AWS 리소스에 대한 권한 오남용 위험 증가
- 감사(Audit) 및 트래픽 추적 어려움

## Diagnosis
서비스 계정(ServiceAccount)에 IRSA 또는 EKS Pod Identity 관련 annotation이 설정되어 있는지 확인합니다.

**Example**
```bash
kubectl get sa --all-namespaces -o jsonpath="{range .items[*]}{.metadata.namespace}{'\t'}{.metadata.name}{'\t'}{.metadata.annotations.eks\.amazonaws\.com/role-arn}{'\t'}{.metadata.annotations.eks\.amazonaws\.com/identity}{'\t'}{.metadata.annotations.eks\.amazonaws\.com/audience}{'\n'}{end}" | grep -v "kube-system"
```

- 출력 결과에서 모든 항목이 비어 있다면, 해당 서비스 계정은 IRSA나 Pod Identity를 사용하지 않는 것입니다.

## Mitigation
**1. IRSA 구성**

IRSA를 구성하려면 IAM Role과 서비스 계정을 연결하는 작업이 필요합니다.

IAM OIDC 공급자 생성
```bash
aws eks describe-cluster --name <cluster-name> --query "cluster.identity.oidc.issuer" --output text
```

IAM Role 생성 (Trust Policy 포함)

Trust Policy는 서비스 계정이 역할을 사용할 수 있도록 허용하는 문서입니다.

서비스 계정에 annotation 추가
```bash
kubectl annotate serviceaccount <service-account-name> \
  -n <namespace> eks.amazonaws.com/role-arn=arn:aws:iam::<account-id>:role/<role-name>
```

**2. EKS Pod Identity 사용**

EKS Pod Identity는 AWS CLI 또는 Console에서 설정 가능하며, 최근 IRSA를 대체할 수 있는 방식으로 자리 잡고 있습니다.

[EKS Pod Identity 공식 문서](https://docs.aws.amazon.com/eks/latest/userguide/pod-identities.html)
[IRSA 공식 문서](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html)
[EKS 보안 모범 사례](https://docs.aws.amazon.com/eks/latest/userguide/security-best-practices.html)