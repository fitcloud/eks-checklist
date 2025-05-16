# SEC-001 EKS 클러스터 API 엔드포인트 접근 제어

## Meaning
Amazon EKS 클러스터의 API 서버는 클러스터 관리, 워크로드 배포 등 핵심 기능을 제공하는 제어 지점입니다.

기본적으로 EKS는 API 엔드포인트를 퍼블릭(공인망) 또는 프라이빗(사설망)으로 구성할 수 있으며, 퍼블릭 접근이 필요한 경우 IP 주소 기반으로 접근을 제한할 수 있습니다.

적절한 접근 제어 없이 퍼블릭으로 열려 있을 경우 보안 위협에 노출될 수 있으며, 특히 외부 공격자가 클러스터를 대상으로 시도하는 무차별 인증 시도 등에 취약해질 수 있습니다.

## Impact
- 보안 위협 노출: 퍼블릭하게 노출된 API는 인증을 우회하거나 취약점을 노린 공격에 노출될 수 있음
- 정보 유출 가능성: 클러스터 설정 정보에 대한 무단 접근 가능
- 컴플라이언스 위반 : 접근 제어 미비로 인해 보안 감사 기준 위반 가능성 존재

## Diagnosis
EKS 클러스터의 API 엔드포인트가 퍼블릭으로 열려 있는지 확인합니다.

**AWS CLI Example**
```bash
aws eks describe-cluster --name <cluster-name> --query "cluster.resourcesVpcConfig"
```
**Result Example**
```json
{
  "endpointPublicAccess": true,
  "publicAccessCidrs": [
    "0.0.0.0/0"
  ],
  "endpointPrivateAccess": false
}
```

- ```endpointPublicAccess```가 ```true```이고

- ```publicAccessCidrs```가 ```0.0.0.0/0```이면 모든 IP에 퍼블릭으로 열려 있는 상태입니다.

## Mitigation
보안을 강화하려면 다음과 같은 조치를 취합니다

**1. 퍼블릭 접근 비활성화 및 프라이빗 접근 활성화**
```bash
aws eks update-cluster-config \
  --region <region> \
  --name <cluster-name> \
  --resources-vpc-config endpointPublicAccess=false,endpointPrivateAccess=true
```

**2. 퍼블릭 접근이 필요한 경우 IP 화이트리스트 설정**
```bash
aws eks update-cluster-config \
  --region <region> \
  --name <cluster-name> \
  --resources-vpc-config endpointPublicAccess=true,publicAccessCidrs="203.0.113.0/24"
```

[EKS 네트워크 구성 가이드](https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html)
[EKS 클러스터 접근 제어 공식 문서](https://docs.aws.amazon.com/eks/latest/userguide/cluster-endpoint.html)