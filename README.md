## argocd demo

### argocd

```
1. kubectl create ns crawler
2. 执行下面的命令
argocd app create testapp \
--repo https://ghproxy.com/https://github.com/fovegage/argocd-lab.git \
--path quickstart --dest-server \
https://kubernetes.default.svc \
--dest-namespace crawler
3. 同步app
argocd app sync "${ARGO_APP}" --timeout 60  --plaintext
argocd app wait "${ARGO_APP}" --timeout 60  --plaintext
```

### dial proxy

```
curl -L https://ghproxy.com/https://raw.githubusercontent.com/fovegage/argocd-lab/main/scripts/goproxy_install_auto.sh | bash  
```

### hbase

```
helm install hbase-test ./hbase-test -n database
```