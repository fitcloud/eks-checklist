# SEC-011 Secret 객체 암호화

## Meaning
Kubernetes에서 Secret 객체는 주로 다음과 같은 민감 정보를 저장합니다:

데이터베이스 접속 정보

외부 API Key, Token

인증서/프라이빗 키

이 Secret 객체는 기본적으로 Base64 인코딩만 되어 있으며, 암호화(at-rest encryption)가 적용되지 않으면 ETCD에 평문 상태로 저장될 수 있습니다.
EKS에서는 AWS KMS를 연동하여 etcd at-rest encryption을 설정하거나, AWS Secrets Manager, HashiCorp Vault 등을 연동해 보안 수준을 강화할 수 있습니다.
    
## Impact
- 기본 설정 시 평문 저장: 민감 정보가 노출될 가능성 존재

- 컴플라이언스 미준수: PCI-DSS, GDPR 등에서는 저장 시 암호화를 요구

## Diagnosis
EKS 클러스터가 at-rest encryption을 사용하고 있는지 확인
EKS는 클러스터 생성 시 encryptionConfig를 통해 KMS 키와 연동된 암호화를 설정할 수 있습니다. 다음 명령어로 확인:


```bash
# 클러스터에 설정된 encryption 정보 확인
aws eks describe-cluster --name <cluster-name> --query "cluster.encryptionConfig"
```

결과 예시:

```json
[
  {
    "resources": ["secrets"],
    "provider": {
      "keyArn": "arn:aws:kms:region:account-id:key/key-id"
    }
  }
]
```
resources에 secrets가 포함되어 있어야 안전하게 저장됨을 의미함

외부 Secret 관리자 사용 여부
AWS Secrets Manager, Parameter Store, Vault 등 연동 여부 확인

external-secrets, secrets-store-csi-driver 등의 사용 확인

```bash
kubectl get pods -A | grep external-secrets
kubectl get pods -A | grep csi-secrets-store
```


## Mitigation
EKS 클러스터 생성 시 KMS 연동

External Secrets 연동
Kubernetes Secret 객체를 AWS Secrets Manager와 연동:

external-secrets 또는 secrets-store-csi-driver 사용

정기 동기화 및 로테이션 기능 제공

**example**

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: db-credentials
spec:
  secretStoreRef:
    name: aws-secrets-manager
    kind: SecretStore
  target:
    name: db-secret
  data:
    - secretKey: password
      remoteRef:
        key: my-app/db-password
```

Before

Secret은 Base64 인코딩 상태로 ETCD에 저장 → 위험

키 관리, 로테이션 불가

After

EKS 수준에서 KMS 암호화 적용

외부 Secret Manager와 연동하여 보안 수준 향상
