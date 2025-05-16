# NET-005 AWS Load Balancer Controller 사용

## Meaning
AWS Load Balancer Controller는 Amazon Web Services(AWS) 환경에서 Kubernetes 클러스터와 함께 사용되는 컨트롤러입니다. 이를 통해 Kubernetes 서비스에 대해 Elastic Load Balancers (ELB)를 자동으로 생성하고 관리할 수 있습니다. 즉, AWS에서 제공하는 로드 밸런서를 Kubernetes와 원활하게 통합하여, 클러스터 외부에서 오는 트래픽을 적절하게 분배하고 관리합니다.

## Impact
- 수동 관리 필요: 로드 밸런서를 수동으로 생성하고 관리해야 하는 포인트가 늘어나며, 설정 오류나 관리의 일관성이 떨어지게 됩니다.
- 유연성 부족: 자동으로 로드 밸런서를 관리하지 않기 때문에 서비스 확장 시 로드 밸런서를 동적으로 조정이 어려움.
- 보안 설정 번거로움: TLS 인증서나 보안 정책 등을 수동으로 설정해야 해서 보안 관리가 어려워질 수 있습니다.

## Diagnosis
AWS Load Balancer Controller가 설치되어 있는지 확인하세요

```bash
kubectl get deployments.apps -A -o custom-columns="NAMESPACE:.metadata.namespace,NAME:.metadata.name" | grep "aws-load-balancer-controller" && echo "PASS" || echo "FAIL"
```

## Mitigation
AWS Load Balancer Controller를 설치하세요

**example**
```bash 
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=my-cluster \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller
```
[AWS Load Balancer Controller란](https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/)
[Helm을 사용한 AWS Load Balancer Controller 설치](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/lbc-helm.html)
[매니페스트를 사용한 AWS Load Balancer Controller 설치](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/lbc-manifest.html)