## argocd demo

```
1. kubectl create ns crawler
2. 执行下面的命令
argocd app create testapp \
--repo https://ghproxy.com/https://github.com/fovegage/argocd-lab.git \
--path quickstart --dest-server \
https://kubernetes.default.svc \
--dest-namespace crawler
```