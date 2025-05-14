# REL-013 Karpenter 기반 노드 생성

## Meaning
Karpenter는 NodeClaim이라는 커스텀 리소스를 통해 노드 생성을 추적하고 관리합니다.
NodeClaim은 Karpenter가 새로운 노드를 프로비저닝할 때 생성되는 리소스로, 다음과 같은 정보를 포함합니다.
- Pod 요구사항을 분석하여 어떤 인스턴스 유형이 필요한지 계산
- 그 조건에 맞는 NodeClass / NodePool을 기반으로 NodeClaim을 생성
- 클라우드 제공자에게 인스턴스를 요청 (launch)
- 생성된 노드를 클러스터에 등록 (register)
- 초기화 완료 시까지 모니터링 (initialize)
NodeClaim이 존재한다는 것은 Karpenter가 워크로드 수요에 반응하여 실제 노드를 생성한 흔적을 의미합니다.

## Impact
NodeClaim이 존재하지 않는 경우 다음과 같은 영향을 받을 수 있습니다.
- Karpenter가 노드를 한 번도 생성하지 않았거나, 설정이 잘못되어 노드 생성에 실패하고 있는 상태일 수 있음
- 리소스 부족 상황에서도 Karpenter가 자동으로 노드를 생성하지 않아 Pod가 Pending 상태로 유지될 수 있음


## Diagnosis
다음 명령어를 통해 NodeClaim 리소스가 존재하는지 확인할 수 있습니다

```bash
kubectl get nodeclaims -A --no-headers 2>/dev/null | grep -q . && echo "Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)가 존재합니다." || echo "Karpenter가 노드를
프로비저닝한 흔적(NodeClaim 리소스)이 존재하지 않습니다."
```

**출력 예시**
정상적인 경우
- Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)가 존재합니다.
비정상적인 경우
- Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)이 존재하지 않습니다.


## Mitigation
다음을 확인하여 문제를 해결합니다.
- Karpenter가 설치되어 있는지 확인
- Karpenter 로그 확인
- CRD(NodeClaim) 리소스가 실제로 존재하는지 확인

[AWS 공식 문서 - Karpenter](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/karpenter.html)
[Karpenter 공식 문서 - Nodeclaim](https://karpenter.sh/docs/concepts/nodeclaims)
[Karpenter 공식 문서 - 설치](https://karpenter.sh/docs/getting-started)
