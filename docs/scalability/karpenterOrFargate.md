# SCL-002 Karpenter 전용 노드 그룹 혹은 Fargate 사용

## Meaning
Karpenter가 관리하는 노드에 Karpenter를 실행하지 않고 최소 하나 이상의 워커 노드가 있는 소규모 전용 노드 그룹 사용 하여 Kapenter를 설치하거나 'Karpenter' 네임스페이스 대한 Fargate Profile을 생성하여 EKS Fagate에서 Karpenter를 실행합니다.

## Impact
- Karpenter가 자기 스스로 종료시킬 수 있습니다
    - Karpenter는 클러스터의 워크로드 요구에따라 노드를 자동으로 확장 및 축소를 진행하게 되며 경우에 따라 Karpenter 자신이 실행중인 노드를 종료 시킬 수 잇습니다.
- 불안정한 오토스케일링 동작
    - Karpenter가 실행중인 노드를 예기치 않게 스케일 다운 또는 교체할 경우 클러스터의 노드 풀 관리가 불안정해집니다.
- 가용성 저하
    - Karpenter 다운되거나 순간 존재 하지않게되면 워크로드 수요에 따라 노드를 생성하지 못하여 서비스 중단으로 이어질 수 있습니다

## Diagnosis
Karpenter 전용 노드 그룹을 확인합니다
```bash

```

Fargate 사용을 확인합니다

```bash

```

## Mitigation
이 기능에 문제가 발생했을 때 적용할 수 있는 완화책을 설명합니다.
