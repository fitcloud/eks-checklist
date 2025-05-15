# SEC-014 읽기 전용 파일시스템 사용

## Meaning
컨테이너에 readOnlyRootFilesystem: true를 설정하면, 루트 파일 시스템을 읽기 전용으로 마운트하여 악의적인 공격 또는 의도치 않은 파일 시스템 변경을 방지할 수 있습니다.
이는 컨테이너 보안을 강화하고, 실행 환경의 불변성(immutability)을 높이는 중요한 보안 수단입니다

## Impact
- 보안 취약점 완화: 악성코드나 취약점을 악용해 루트 파일 시스템을 변경하는 행위를 방지합니다.
- 불변성 확보: readOnlyRootFilesystem 설정하면, 루트 디렉터리에 쓰기가 불가능 해지기 때문에 시스템 전반에 영향을 미치는 변경이 원천 차단됨

## Diagnosis
읽기 전용 파일시스템을 사용중인지 확인하세요
```bash
kubectl get pods --all-namespaces -o json | jq '
  .items[] |
  select(.metadata.namespace != "kube-system") |
  {
    namespace: .metadata.namespace,
    pod: .metadata.name,
    containers: [
      .spec.containers[] |
      select(
        (.securityContext.readOnlyRootFilesystem != true)
      ) |
      .name
    ]
  } | select(.containers | length > 0)'
```
## Mitigation
읽기 전용 파일시스템 설정을 추가하세요

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 3000
    fsGroup: 2000
    supplementalGroups: [4000]
  volumes:
  - name: sec-ctx-vol
    emptyDir: {}
  containers:
  - name: sec-ctx-demo
    image: busybox:1.28
    command: [ "sh", "-c", "sleep 1h" ]
    volumeMounts:
    - name: sec-ctx-vol
      mountPath: /data/demo
securityContext:
  readOnlyRootFilesystem: true
```

[Kubernetes 공식 문서 - Security Context](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
[Kubernetes 공식 문서 - Immutable container filesystems](https://kubernetes.io/blog/2021/10/05/nsa-cisa-kubernetes-hardening-guidance/#immutable-container-filesystems)
[Amazon EKS 환경에서 Pod Security Standard 구현하기](https://aws.amazon.com/ko/blogs/tech/implementing-pod-security-standards-in-amazon-eks/)