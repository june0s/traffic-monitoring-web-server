# Traffic Monitoring web server (DEMO)
 

- istio + kiali 를 보다 쉽게 사용하기 위해 만들어 본 데모
- K8s 의 client go 를 사용하기 위해 go 언어로 작성

## 소스 구조

```
.
├── README.md
├── go.mod
├── go.sum
├── main.go
└── pkg
    ├── restapi
    │   ├── k8sRestApi.go
    │   └── response.go
    └── utils
        └── k8sUtils.go
```

* pkg/restapi : 클라이언트로 부터 restapi 요청 시 응답하는 controller 소스
* pkg/utils : K8s api 와 통신하여 네임스페이스 리스트 등 정보 가져오는 소스 

## 사전 준비 (로컬 테스트)

1. kubeconfig 파일 준비 - K8s 개발 환경의 kubeconfig 을 준비
2. kubeconfig 파일 설정 - pkg/utils/k8sUtils.go 파일의 ctxValue 에 1번 kubeconfig 파일 경로 설정 
3. kiali 포트 포워딩 - kiali 서버를 로컬에서 접속할 수 있게 포트 포워딩 해둔다.
```
$ kubectl -n istio-system port-forward svc/kiali 20001
```

## 실행 방법

1. main.go 소스가 있는 위치에서 아래 명령어 실행
```
$ go run main.go
```
