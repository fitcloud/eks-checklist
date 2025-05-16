# NET-004 사용 사례에 맞는 로드밸런서 사용(ALB or NLB)

## Meaning
AWS에서 제공하는 **Application Load Balancer(ALB)**와 **Network Load Balancer(NLB)**는 각기 다른 특성과 용도를 가지고 있습니다. 사용 사례에 맞게 적절한 로드밸런서를 선택하여, 효율적이고 안정적인 트래픽 분배를 할 수 있도록 합니다.
- ALB: HTTP/HTTPS 기반 애플리케이션에 적합, URL 경로, 호스트 이름 기반의 라우팅 지원

- NLB: TCP/UDP 트래픽 처리에 적합, 고성능 및 저지연이 필요한 서비스에 유리
## Impact
- ALB vs NLB 부적절 사용: 잘못된 로드밸런서를 선택하면 성능 저하, 비용 증가 및 애플리케이션 장애가 발생할 수 있음

- 잘못된 라우팅: 애플리케이션의 요구 사항에 맞지 않는 로드밸런서를 사용하면 트래픽이 제대로 라우팅되지 않을 수 있음

## Diagnosis
애플리케이션의 트래픽 유형과 요구 사항을 기준으로 ALB 또는 NLB를 선택합니다.

HTTP/HTTPS 트래픽이 필요한 경우 ALB를 선택합니다.

TCP/UDP와 같은 고성능 트래픽 처리나 지연 최소화가 중요한 경우 NLB를 선택합니다.

## Mitigation
적절한 로드밸런서를 선택하여 설정합니다.

**ALB 설정 예시 (애플리케이션 트래픽에 적합)**

```bash
aws elbv2 create-load-balancer --name my-alb --subnets subnet-12345678 --security-groups sg-12345678 --scheme internet-facing --load-balancer-type application
```

**NLB 설정 예시 (고성능 TCP/UDP 트래픽에 적합)**

```bash
aws elbv2 create-load-balancer --name my-nlb --subnets subnet-12345678 --security-groups sg-12345678 --scheme internet-facing --load-balancer-type network
```

[AWS ALB NLB 사용 사례](https://aws.amazon.com/ko/compare/the-difference-between-the-difference-between-application-network-and-gateway-load-balancing/)