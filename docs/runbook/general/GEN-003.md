# GEN-003 컨테이너 이미지 태그에 latest 미사용

## Meaning
컨테이너 이미지를 사용할 때 latest 태그를 사용하는 것은 일반적으로 권장되지 않습니다.
latest는 실제 어떤 버전의 이미지인지 명확하지 않으며, 시간이 지남에 따라 이미지 내용이 바뀔 수 있어 예측 불가능한 배포나 동작을 유발할 수 있습니다.
지속적인 운영 및 디버깅, 롤백 등을 고려할 때 명시적인 버전 태그(예: v1.0.3)를 사용하는 것이 바람직합니다.

## Impact
- 불안정한 배포: ```latest``` 태그가 변경될 경우, 동일한 배포라도 실행 환경이 달라질 수 있습니다.
- 디버깅 어려움: 문제가 발생했을 때 정확한 이미지 버전을 알 수 없어 원인 분석이 어려워집니다.
- 롤백 어려움: 이전 상태로 쉽게 되돌릴 수 없어 운영 안정성 저하 가능성.

## Diagnosis
클러스터 내 모든 Pod의 컨테이너 이미지에서 ```latest``` 태그가 사용되었는지 확인합니다.
```bash
kubectl get pods --all-namespaces -o jsonpath='{range .items[*]}{.metadata.namespace}{"\t"}{.metadata.name}{"\t"}{range .spec.containers[*]}{.name}{"\t"}{.image}{"\n"}{end}{end}' | grep ':latest'

<!--
kubectl get pods -A -o custom-columns=NAMESPACE:.metadata.namespace,NAME:.metadata.name,IMAGE:.spec.containers[*].image | awk 'NR==1 || $3 ~ /:latest/'
-->
```
## Mitigation
컨테이너 이미지에서 ```latest``` 태그 대신 명시적인 버전 태그를 사용하도록 합니다.

**example**

**Before**
```yaml
image: nginx:latest
```

**After**
```yaml
image: nginx:1.25.2
```

**이미지 버전 확인 방법**
DockerHub 또는 이미지 저장소에서 해당 이미지의 사용 가능한 태그를 확인한 후, 안정적인 버전을 선택하여 사용합니다.

```bash
# DockerHub에서 nginx 태그 확인 예시
curl -s https://registry.hub.docker.com/v1/repositories/nginx/tags | jq '.[].name'
```

[컨테이너 이미지 관리 가이드](https://kubernetes.io/ko/docs/concepts/containers/images/#updating-images)