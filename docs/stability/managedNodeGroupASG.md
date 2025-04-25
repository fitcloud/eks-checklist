# **오토스케일링_그룹_기반_관리형_노드_그룹_생성**

## Meaning
Amazon EKS에서 관리형 노드 그룹(Managed Node Group) 은 AWS Auto Scaling Group(ASG)을 통해 관리되며, 클러스터의 컴퓨팅 자원을 유연하게 확장하거나 축소할 수 있도록 되어 있습니다.
자동 확장을 올바르게 구성하려면 ASG의 minSize 값이 maxSize보다 작아야 하며, 이를 통해 워크로드 수요에 따라 동적으로 노드를 추가하거나 제거할 수 있습니다.
만약 이 구성이 올바르지 않거나 누락되어 있다면, 클러스터는 확장성과 복원력을 갖추지 못하게 됩니다.


## Impact
- 스케일링 불가: minSize ≥ maxSize일 경우, 클러스터는 현재 설정된 용량 외에 확장될 수 없어 수요에 따른 유연한 대응이 불가능합니다.

- 과도한 고정 비용: 모든 노드를 고정된 수량으로 유지하게 되면, 이는 리소스 과잉으로 이어져 불필요한 비용이 발생할 수 있습니다.

- 장애 대응력 부족: Auto Scaling이 동작하지 않으면, 장애 복구 시 신규 노드 생성성이 불가능해 전체 서비스에 영향을 줄 수 있습니다.

## Diagnosis
다음 명령어를 실행하면 현재 클러스터에 ASG 기반 관리형 노드 그룹이 존재하는지, 그리고 해당 그룹이 올바르게 스케일링 구성(minSize < maxSize) 되어 있는지를 확인할 수 있습니다.
CLUSTER_NAME, AWS_REGION, AWS_PROFILE은 환경변수로 사전에 등록이 되어있어야 합니다.

```bash
# 변수 값 설정정
CLUSTER_NAME=<CLUSTER_NAME>
AWS_REGION=<AWS_REGION>
AWS_PROFILE=<AWS_PROFILE>

NODEGROUPS=$(kubectl get nodes -o json | jq -r '.items[].metadata.labels["eks.amazonaws.com/nodegroup"]' | sort -u | grep -v null)

if [ -z "$NODEGROUPS" ]; then
  echo "관리형 노드그룹이 존재하지 않습니다."
  exit 0
fi

echo "$NODEGROUPS" | while read ng; do
  asg=$(aws eks describe-nodegroup --cluster-name "$CLUSTER_NAME" --nodegroup-name "$ng" --region "$AWS_REGION" --profile "$AWS_PROFILE" | jq -r '.nodegroup.resources.autoScalingGroups[0].name')
  if [ "$asg" == "null" ] || [ -z "$asg" ]; then
    echo "$ng (ASG 없음)"
    continue
  fi
  aws autoscaling describe-auto-scaling-groups --auto-scaling-group-names "$asg" --region "$AWS_REGION" --profile "$AWS_PROFILE" \
  | jq -r --arg ng "$ng" --arg asg "$asg" '
    .AutoScalingGroups[0] as $g
    | if ($g.MinSize < $g.MaxSize) then
        "Nodegroup: \($ng) | ASG: \($asg) (minSize: \($g.MinSize), maxSize: \($g.MaxSize))"
      else
        "\($ng) (minSize ≥ maxSize)"
      end'
done
```
### 출력 예시시
- 관리형 노드그룹이 존재하지 않습니다. → 관리형 노드가 클러스터에 없음

- ng-a (ASG 없음) → 노드그룹은 있으나 ASG가 할당되어 있지 않음

- ng-b (minSize ≥ maxSize) → 자동 확장이 비활성화 상태

- Nodegroup: ng-c | ASG: asg-c (minSize: 1, maxSize: 4) → 정상 구성


## Mitigation

문제가 되는 노드그룹에 대해 아래와 같이 수정합니다.

### AWS CLI 또는 콘솔에서 해당 Auto Scaling Group의 minSize / maxSize 값을 조정해야 합니다
example
```bash
aws autoscaling update-auto-scaling-group \
  --auto-scaling-group-name <ASG_NAME> \
  --min-size 1 \
  --max-size 5 \
  --region <REGION>
```
### [AWS Autoscaling 공식문서](https://docs.aws.amazon.com/cli/latest/reference/autoscaling/update-auto-scaling-group.html)
### [AWS Github.io Cluster Autoscaler](https://aws.github.io/aws-eks-best-practices/ko/cluster-autoscaling/) 