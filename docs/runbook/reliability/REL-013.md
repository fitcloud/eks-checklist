# REL-013 Karpenter 기반 노드 생성

## Meaning
Karpenter는 NodeClaim을 통해 노드를 자동으로 생성하고 관리합니다. NodeClaim은 Karpenter가 노드를 프로비저닝한 증거을 의미합니다.

## Impact
- NodeClaim이 없으면 Karpenter가 노드를 생성하지 않거나 설정 오류로 생성 실패 가능
- 리소스 부족 시 자동으로 노드가 생성되지 않아 Pod가 Pending 상태로 유지될 수 있습니다

## Diagnosis
다음 명령어로 NodeClaim 리소스를 확인합니다

```bash
kubectl get nodeclaims -A --no-headers 2>/dev/null | grep -q . && echo "Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)가 존재합니다." || echo "Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)이 존재하지 않습니다."
```

**출력 예시**
정상적인 경우
- Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)가 존재합니다.
비정상적인 경우
- Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)이 존재하지 않습니다.

## Mitigation
Karpenter가 정상적으로 설치되었는지 확인하고, 설치가 안 되어 있으면 설치합니다.

NodeClaim 리소스가 생성되지 않는 경우 Karpenter의 설정과 NodeClass, NodePool이 올바르게 구성되었는지 점검합니다.

[AWS 공식 문서 - Karpenter](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/karpenter.html)
[Karpenter 공식 문서 - Nodeclaim](https://karpenter.sh/docs/concepts/nodeclaims)
[Karpenter 공식 문서 - 설치](https://karpenter.sh/docs/getting-started)
