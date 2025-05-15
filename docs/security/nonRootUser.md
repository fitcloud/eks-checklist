# SEC-005 루트 유저가 아닌 유저로 컨테이너 실행

## Meaning
이 점검은 Kubernetes 컨테이너가 root 사용자(UID 0) 로 실행되고 있는지 확인하는 보안 검사입니다.
Kubernetes에서 securityContext.runAsUser 필드를 통해 컨테이너를 일반 사용자 권한으로 실행할 수 있으며, 해당 필드가 명시적으로 설정되어 있고 UID가 0이 아닐 때만 안전한 상태로 판단됩니다.
또한 Windows 컨테이너의 경우, runAsUserName이 "Administrator"로 설정되면 루트 권한과 유사한 위험으로 간주합니다.
컨테이너가 runAsUser: 0(root) 또는 미지정 상태면 보안 취약. 반드시 일반 사용자 권한(예: UID 1000)으로 실행해야 안전합니다.

## Impact
- 루트 권한 컨테이너는 취약점 발생 시 호스트 시스템까지 침해 가능
- Pod 간 보안 격리 실패 가능성 증가
- 보안사고 확산 가능성 증가

## Diagnosis
아래 명령어를 사용하면 클러스터 내의 컨테이너 중에서 다음 조건에 해당하는 경우를 탐지할 수 있습니다.
- runAsUser가 0으로 명시된 경우 (명시적으로 root 권한 실행)
- securityContext 또는 runAsUser가 아예 명시되지 않은 경우 (기본적으로 root 권한 실행 가능)
- Windows 컨테이너에서 runAsUserName이 "Administrator"로 설정된 경우

**command example**
```bash
kubectl get pods -A -o json | jq -r '
.items[] | . as $p | $p.spec.containers[]
| select((.securityContext.runAsUser // 0) == 0)
| "NS:\($p.metadata.namespace) Pod:\($p.metadata.name) Ctr:\(.name) (root 실행)"
'
```
- 이 명령은 보안 상 취약한 컨테이너만 필터링하여 표시하므로, 출력이 없다면 모든 컨테이너가 적절한 사용자 권한으로 실행되고 있는 것입니다.

**출력 example**
```bash
Namespace: default | Pod: nginx | Container: nginx (RunAsUser 미설정, root로 실행 가능성 존재)
Namespace: test | Pod: backend | Container: app (명시적 root 계정 실행)
Namespace: winspace | Pod: winpod | Container: winapp (Windows Administrator 실행)
```

- RunAsUser 미설정: securityContext 또는 runAsUser가 지정되지 않아 기본적으로 root 권한으로 실행될 가능성이 있습니다.
- 명시적 root 계정 실행: runAsUser: 0으로 명시되어 있어 의도적으로 루트로 실행되도록 설정된 컨테이너입니다.
- Windows Administrator 실행: Windows 기반 컨테이너에서 runAsUserName이 "Administrator"로 설정된 경우로, 보안상 root 실행과 동일한 위험을 가집니다.

## Mitigation
컨테이너에 securityContext.runAsUser를 명시하고 UID 0을 피하세요. 루트 권한을 회피하는 것은 기본 보안 원칙 중 하나입니다.

**Non-Root 사용자로 실행 설정 example**
```yaml
securityContext:
  runAsUser: 1000
  allowPrivilegeEscalation: false
```

[Pod-security](https://kubernetes.io/ko/docs/concepts/security/pod-security-standards/)