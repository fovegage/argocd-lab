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
kubectl create -f https://ghproxy.com/https://raw.githubusercontent.com/pingcap/tidb-operator/v1.3.9/manifests/crd.yaml -n storage


helm install hbase-test ./hbase-test -n database

# upgrade  install
helm upgrade tidb-operator ./tidb-operator-1.4.0-beta.3/charts/tidb-operator   -f ./tidb-operator-1.4.0-beta.3/charts/tidb-operator/values.yaml --namespace=storage

//
    kubectl get pods --namespace storage -l app.kubernetes.io/instance=tidb-operator


helm install tidb-cluster  ./tidb-operator-1.4.0-beta.3/charts/tidb-cluster -f ./tidb-operator-1.4.0-beta.3/charts/tidb-cluster/values.yaml  --namespace storage

//
1. Watch tidb-cluster up and running
     watch kubectl get pods --namespace storage -l app.kubernetes.io/instance=tidb-cluster -o wide
2. List services in the tidb-cluster
     kubectl get services --namespace storage -l app.kubernetes.io/instance=tidb-cluster

Cluster access
* Access tidb-cluster using the MySQL client
    kubectl port-forward -n storage svc/tidb-cluster-tidb 4000:4000 &
    mysql -h 127.0.0.1 -P 4000 -u root -D test
  Set a password for your user
    SET PASSWORD FOR 'root'@'%' = 'uQbMhhohHU'; FLUSH PRIVILEGES;
* View monitor dashboard for TiDB cluster
   kubectl port-forward -n storage svc/tidb-cluster-grafana 3000:3000
   Open browser at http://localhost:3000. The default username and password is admin/admin.
   If you are running this from a remote machine, you must specify the server's external IP address.
```

```
# 配置
kubectl port-forward -n storage svc/tidb-cluster-tidb 4000:4000 


https://asktug.com/t/topic/34231

echo "fs.file-max = 1000000" >> /etc/sysctl.conf
sysctl -p 

vi /usr/lib/systemd/system/docker.service
LimitNOFILE=1000000
LimitNPROC=1000000

# 各个机器
systemctl daemon-reload
systemctl restart docker

kubectl port-forward service/tidb-cluster-tidb 31151 -n storage
```

```
# google 常用的国内镜像
https://www.cnblogs.com/w84422/p/15596883.html
```

```
# metabase
https://juejin.cn/post/7169124983015211022

```

```
kubectl delete -f https://ghproxy.com/https://raw.githubusercontent.com/pingcap/tidb-operator/v1.3.9/manifests/crd.yaml -n storage

helm upgrade tidb-operator ./tidb-1.3.9/charts/tidb-operator   -f ./tidb-1.3.9/charts/tidb-operator/values.yaml --namespace=storage
helm upgrade tidb-cluster  ./tidb-1.3.9/charts/tidb-cluster -f ./tidb-1.3.9/charts/tidb-cluster/values.yaml  --namespace=storage

helm upgrade tidb-operator ./tidb-operator-1.4.0-beta.3/charts/tidb-operator   -f ./tidb-operator-1.4.0-beta.3/charts/tidb-operator/values.yaml --namespace=storage
helm upgrade tidb-operator ./tidb-operator-1.4.0-beta.3/charts/tidb-operator   -f ./tidb-operator-1.4.0-beta.3/charts/tidb-operator/values.yaml --namespace=storage

```

```
helm list -n storage
helm uninstall tidb-cluster -n storage
```

```
kubectl create -f https://ghproxy.com/https://raw.githubusercontent.com/pingcap/tidb-operator/v1.3.9/manifests/crd.yaml

```