# **데이터 플레인 노드에 필수로 필요한 IAM 권한만 부여**

## Meaning
데이터 플레인 노드(워커 노드)에는 최소한의 필수 IAM 권한만 부여되어야 합니다.  
이 검사는 노드에 연결된 IAM 역할이 **허용된 정책** (AmazonEC2ContainerRegistryReadOnly, AmazonEKS_CNI_Policy, AmazonEKSWorkerNodePolicy) 외의 다른 권한을 포함하고 있지 않은지를 확인합니다.  

## Impact
- 노드에 불필요한 권한이 부여되면, 노드가 침해될 경우 그 노드에 할당된 모든 권한이 악용될 수 있습니다.
- 불필요하게 부여된 추가 권한을 통해 클러스터 내부 뿐만 아니라 AWS 계정 내 다른 리소스에 접근할 수 있습니다.
- 최소 권한 원칙(Least Privilege Principle)을 위반하면, 감사 및 규정 준수 측면에서도 문제가 발생합니다.

## Diagnosis

```bash
node_ips=($(kubectl get nodes -o custom-columns="NAME:.metadata.name,IP:.metadata.annotations.alpha\.kubernetes\.io/provided-node-ip" | tail -n +2 | awk '{print $2}'))
echo "Node IPs: ${node_ips[@]}"

allowed_regex='^(AmazonEC2ContainerRegistryReadOnly|AmazonEKS_CNI_Policy|AmazonEKSWorkerNodePolicy)$'

# 배열에 있는 각 노드 IP에 대해 반복 실행
for node_ip in "${node_ips[@]}"; do
  echo "Processing node IP: ${node_ip}"

  profile_arn=$(aws ec2 describe-instances \
    --filters "Name=private-ip-address,Values=${node_ip}" \
    --query "Reservations[].Instances[].IamInstanceProfile.Arn" \
    --output json | jq -r '.[0]')
  
  if [ -z "$profile_arn" ] || [ "$profile_arn" == "null" ]; then
    echo "  Node IP ${node_ip}: No IAM instance profile found."
    echo "-----------------------------------"
    continue
  fi
  echo "  Instance Profile ARN: ${profile_arn}"

  instance_profile=$(basename "$profile_arn")
  echo "  Instance Profile: ${instance_profile}"

  role_name=$(aws iam get-instance-profile \
    --instance-profile-name "${instance_profile}" \
    --query "InstanceProfile.Roles[].RoleName" \
    --output json | jq -r '.[0]')
  
  if [ -z "$role_name" ] || [ "$role_name" == "null" ]; then
    echo "  Node IP ${node_ip}: No IAM role found in instance profile ${instance_profile}."
    echo "-----------------------------------"
    continue
  fi
  echo "  IAM Role: ${role_name}"

  attached_policies=$(aws iam list-attached-role-policies \
    --role-name "${role_name}" \
    --output json --query "AttachedPolicies[].PolicyName" | jq -r '.[]')
  
  echo "  Attached policies:"
  echo "${attached_policies}"
  
  # 허용된 정책을 제외한 나머지(허용되지 않은 정책) 필터링
  unauthorized_policies=$(echo "${attached_policies}" | grep -v -x -E "${allowed_regex}")
  
  if [ -n "${unauthorized_policies}" ]; then
    echo "  Node IP: ${node_ip} (Role: ${role_name}) -> Unauthorized policies:"
    echo "${unauthorized_policies}"
  else
    echo "  Node IP: ${node_ip} (Role: ${role_name}) -> Only allowed policies are attached."
  fi
  echo "-----------------------------------"
done

```

## Mitigation
데이터 플레인 노드에 반드시 필요한 권한만 포함하도록 IAM 역할 및 인스턴스 프로파일을 확인합니다.

불필요한 정책 제거
```bash
aws iam detach-role-policy --role-name <MyNodeRole> --policy-arn <arn:aws:iam::aws:policy/UnwantedPolicyName>
```
