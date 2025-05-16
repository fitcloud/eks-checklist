# REL-011 오토스케일링 그룹 기반 관리형 노드 그룹 생성

## Meaning
EKS에서 관리형 노드 그룹은 Auto Scaling Group(ASG)으로 유동적으로 확장되며, minSize가 maxSize보다 작아야만 확장 가능합니다. 잘못 설정된 경우 확장되지 않거나 비용이 불필요하게 증가할 수 있습니다.

## Impact
- 스케일링 불가: minSize ≥ maxSize일 경우 자동 확장이 불가능.
- 비용 증가: 고정된 노드 수로 과도한 비용 발생.
- 장애 대응 불가: Auto Scaling 미작동 시 노드 추가 불가.

## Diagnosis
ASG의 minSize와 maxSize가 올바르게 설정되었는지 확인합니다.

```bash
kubectl get nodes -o json | jq -r '.items[].metadata.labels["eks.amazonaws.com/nodegroup"]' | sort -u | grep -v null | xargs -I {} aws autoscaling describe-auto-scaling-groups --auto-scaling-group-names {} --region <AWS_REGION> --profile <AWS_PROFILE> --query 'AutoScalingGroups[0].{Name:AutoScalingGroupName, MinSize:MinSize, MaxSize:MaxSize}'
```
## Mitigation

ASG의 minSize와 maxSize 값을 올바르게 설정합니다.

```bash
aws autoscaling update-auto-scaling-group --auto-scaling-group-name <ASG_NAME> --min-size 1 --max-size 5 --region <REGION>
```
[AWS Autoscaling 공식문서](https://docs.aws.amazon.com/cli/latest/reference/autoscaling/update-auto-scaling-group.html)
[AWS Github.io Cluster Autoscaler](https://aws.github.io/aws-eks-best-practices/ko/cluster-autoscaling/) 