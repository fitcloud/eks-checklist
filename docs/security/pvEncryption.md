# SEC-010 PV 암호화

## Meaning
EKS에서 Pod가 사용하는 영구 저장소는 보통 **EBS 볼륨을 통한 PersistentVolume(PV)**으로 구성됩니다. 이러한 저장소는 데이터가 평문 상태로 저장되면 유출 시 치명적인 보안 위협이 될 수 있으므로, 반드시 암호화를 적용해야 합니다.

**EBS 암호화(EBS Encryption)**를 사용하면, AWS에서 관리하거나 사용자 지정한 KMS 키를 통해 암호화된 볼륨이 자동 생성되며, 해당 볼륨을 사용하는 PVC도 암호화됩니다.

## Impact
- 데이터 유출 시 위험: 볼륨 탈취 또는 Snapshot 복제를 통한 평문 데이터 유출 가능성
- 복구 시 데이터 노출: Snapshot 백업/복구 시에도 평문으로 노출 가능

## Diagnosis
EKS에서 사용 중인 PersistentVolume(PV)과 PVC(PersistentVolumeClaim)를 조회하고, 해당 PV가 연결된 EBS 볼륨의 암호화 여부를 확인합니다.


```bash
# PVC와 연결된 PV 확인
kubectl get pvc -A -o wide

# PV에서 VolumeID(EBS ID) 확인
kubectl get pv -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.awsElasticBlockStore.volumeID}{"\n"}{end}'

# AWS CLI를 통한 EBS 암호화 여부 확인
aws ec2 describe-volumes --volume-ids vol-xxxxxxxx | jq '.Volumes[].Encrypted'
Encrypted: true 가 나와야 암호화된 상태입니다.
```

사용중인 스토리지 클래스에 다음과 같은 설정이 있는지 확인할 수 있습니다

```yaml
parameters:
  encrypted: "true"
  kmsKeyId: <optional-custom-kms-key-id>
```

## Mitigation
다음과 같은 방식으로 EKS에서 사용할 PV에 암호화를 적용합니다.

스토리지 클래스(StorageClass) 수정
기본 StorageClass 또는 사용자 정의 StorageClass에서 EBS 암호화 옵션을 활성화합니다.

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: encrypted-gp3
provisioner: ebs.csi.aws.com
parameters:
  type: gp3
  encrypted: "true"
  kmsKeyId: arn:aws:kms:ap-northeast-2:123456789012:key/abcde-1234-5678-xyz
```

encrypted: "true"로 설정

KMS 키를 명시적으로 지정 가능 (기본 AWS 관리 키 사용도 가능)

**Before**

StorageClass에 encrypted 설정 없음 → 평문 EBS 사용

**After**

모든 PVC는 암호화된 EBS를 기반으로 생성, KMS 키 관리로 보안 수준 향상

[AWS EKS PV 암호화](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/data-encryption-and-secrets-management.html)