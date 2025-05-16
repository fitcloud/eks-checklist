# SEC-002 클러스터 접근 제어(Access entries, aws-auth 컨피그맵)

## Meaning
클러스터 접근 제어는 Kubernetes 클러스터에 누가 접근할 수 있는지와 어떤 권한을 가지는지를 제어하는 핵심 보안 기능입니다. 이 점검은 클러스터에 대해 Access Entries 또는 aws-auth ConfigMap이 적절히 설정되어 있는지를 확인합니다.

## Impact
- 적절한 권한 부여: 필요한 사용자가 권한을 받지 못하면 운영/관리 업무에 차질이 발생할 수 있으며, aws-auth 경우 수동으로 잘못 편집하면 클러스터에서 모든 접근 권한이 차단되는 위험도 존재합니다
- 보안 사고: 클러스터 접근 권한이 잘못 구성될 경우, 비인가 사용자가 클러스터에 접근할 수 있어 보안 사고로 이어질 수 있습니다.

## Diagnosis
```bash
# EKS Configmap aws-auth 확인
kubectl get configmap aws-auth -n kube-system -o yaml
# EKS Cluster Access Entry 확인
aws eks list-access-entries --cluster-name <cluster-name>
```
## Mitigation

- EKS 클러스터에 접근해야 하는 사용자/역할만 aws-auth 또는 Access Entry에 명시되어 있는지 수동으로 확인하세요

- 향후 AWS EKS의 특정 Kubenetes 버전에서는 지원되는 인증 소스에서 aws-auth Configmap이 제거될 예정이므로, Access Entry 방식으로 전환하여 정책 기반 관리와 콘솔 UI를 통한 편리한 권한 제어을 권장합니다.

[EKS 보안 ID 및 액세스 관리](https://docs.aws.amazon.com/eks/latest/best-practices/identity-and-access-management.html)
[간소화된 Amazon EKS 액세스 관리 제어](https://aws.amazon.com/ko/blogs/tech/a-deep-dive-into-simplified-amazon-eks-access-management-controls/)