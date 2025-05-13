# REL-005-Probe(Startup, Readiness, Liveness) 적용

## Meaning
Kubernetes에서는 애플리케이션의 상태를 주기적으로 확인하기 위해 3가지 Probe를 제공합니다.

- startupProbe: 컨테이너 내의 애플리케이션이 시작되었는지를 나타냅니다. 스타트업 프로브(startup probe)가 주어진 경우, 성공할 때까지 다른 나머지 프로브는 활성화되지 않습니다. 만약 스타트업 프로브가 실패하면, kubelet이 컨테이너를 죽이고, 컨테이너는 재시작 정책에 따라 처리됩니다. 컨테이너에 스타트업 프로브가 없는 경우, 기본 상태는 Success 입니다.

- readinessProbe: 컨테이너가 요청을 처리할 준비가 되었는지 여부를 나타냅니다. 만약 준비성 프로브(readiness probe)가 실패한다면, 엔드포인트 컨트롤러는 파드에 연관된 모든 서비스들의 엔드포인트에서 파드의 IP주소를 제거합니다. 준비성 프로브의 초기 지연 이전의 기본 상태는 Failure 입니다. 만약 컨테이너가 준비성 프로브를 지원하지 않는다면, 기본 상태는 Success 입니다.

- livenessProbe: 컨테이너가 동작 중인지 여부를 나타냅니다. 만약 활성 프로브(liveness probe)에 실패한다면, kubelet은 컨테이너를 죽이고, 해당 컨테이너는 재시작 정책의 대상이 됩니다. 만약 컨테이너가 활성 프로브를 제공하지 않는 경우, 기본 상태는 Success 입니다.

이 Probe들을 활용하면 애플리케이션의 신뢰성과 자동 복구 능력을 높일 수 있으며, 실제 서비스 중단 없이 장애 대응이 가능해집니다.

## Impact
- 느린 부팅 앱 안정화: startupProbe를 사용하면 앱이 준비될 때까지 livenessProbe가 재시작하지 않도록 보호할 수 있습니다.

- 서비스 중단 방지: readinessProbe는 준비되지 않은 Pod에 트래픽이 전달되지 않도록 하여 장애 전파를 방지합니다.

- 빠른 장애 감지 및 복구: livenessProbe를 통해 컨테이너가 비정상 상태일 경우 빠르게 감지하고 재시작할 수 있습니다.


## Diagnosis
아래 명령어를 실행하면, kube-system을 제외한 모든 파드 중 probe 미설정 컨테이너를 확인할 수 있습니다.

```bash
kubectl get pods --all-namespaces -o json | jq -r '
  .items[]
  | select(.metadata.namespace != "kube-system")
  | . as $pod
  | $pod.spec.containers[]
  | {
      namespace: $pod.metadata.namespace,
      pod: $pod.metadata.name,
      container: .name,
      missing: [
        (if has("startupProbe") | not then "startupProbe" else empty end),
        (if has("livenessProbe") | not then "livenessProbe" else empty end),
        (if has("readinessProbe") | not then "readinessProbe" else empty end)
      ]
    }
  | select(.missing | length > 0)
  | "Namespace: \(.namespace) | Pod: \(.pod) | Container: \(.container) (미설정: \(.missing))"
'
```

## Mitigation
파드 정의에서 모든 컨테이너에 대해 아래와 같은 Probe 설정을 추가하세요.
다음 예시는 livenessProbe의 사용 예시 입니다.
example
```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    test: liveness
  name: liveness-exec
spec:
  containers:
  - name: liveness
    image: registry.k8s.io/busybox:1.27.2
    args:
    - /bin/sh
    - -c
    - touch /tmp/healthy; sleep 30; rm -f /tmp/healthy; sleep 600
    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/healthy
      initialDelaySeconds: 5
      periodSeconds: 5

```
### [프로브 사용법](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
