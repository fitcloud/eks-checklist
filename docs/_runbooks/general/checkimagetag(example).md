---
title: "checkimagetag"
layout: single
categories: ["general"]
permalink: /runbooks/general/checkimagetag(example)/
---

# checkimagetag

## 문제 정의
Kubernetes Horizontal Pod Autoscaler (HPA)가 최대치에 도달했을 때 발생하는 문제를 해결하는 가이드입니다.

## 해결 방법
1. `kubectl get hpa` 명령어를 실행하여 현재 HPA 상태를 확인합니다.
2. 리소스 제한을 초과하는지 체크합니다.
3. 필요하다면 `spec.maxReplicas` 값을 조정합니다.

## 추가 정보
- [공식 문서](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
