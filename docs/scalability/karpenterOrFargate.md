# SCL-002 Karpenter 전용 노드 그룹 혹은 Fargate 사용

## Meaning
Karpenter가 자신이 실행 중인 노드를 잘못 종료할 위험이 있으므로, 안정적인 운영을 위해 Karpenter 전용 노드 그룹 또는 Fargate를 사용하여 Karpenter 자체는 항상 안정적으로 실행되도록 구성해야 합니다.

## Impact
- Karpenter가 자기 스스로 종료시킬 수 있습니다: Karpenter는 클러스터의 워크로드 요구에따라 노드를 자동으로 확장 및 축소를 진행하게 되며 경우에 따라 Karpenter 자신이 실행중인 노드를 종료 시킬 수 잇습니다.
- 불안정한 오토스케일링 동작: Karpenter가 실행중인 노드를 예기치 않게 스케일 다운 또는 교체할 경우 클러스터의 노드 풀 관리가 불안정해집니다.
- 가용성 저하: Karpenter 다운되거나 순간 존재 하지않게되면 워크로드 수요에 따라 노드를 생성하지 못하여 서비스 중단으로 이어질 수 있습니다

## Diagnosis
Karpenter 전용 노드 그룹을 확인합니다
```bash
kubectl get node -l karpenter.sh/nodepool -o wide
```
Fargate 사용을 확인합니다

```bash
kubectl get node -l eks.amazonaws.com/fargate-profile -o wide
```

## Mitigation
Karpenter가 관리하는 노드에 Karpenter를 실행하지 않고 최소 하나 이상의 워커 노드가 있는 소규모 전용 노드 그룹 사용 하여 Kapenter를 설치하거나 'Karpenter' 네임스페이스 대한 Fargate Profile을 생성하여 EKS Fagate에서 Karpenter를 실행하세요.

**example**

전용 노드 그룹
```bash
eksctl create nodegroup \
  --cluster <클러스터-이름> \
  --name <노드그룹-이름> \
  --node-type <인스턴스-타입> \
  --nodes 2 \
  --managed
```
Fargate

```bash
aws eks create-fargate-profile \
  --cluster-name <your-cluster-name> \
  --fargate-profile-name karpenter-fargate-profile \
  --namespace karpenter \
  --pod-execution-role-arn <your-fargate-pod-execution-role-arn>
```

[Karpenter 시작하기](https://karpenter.sh/docs/getting-started/)