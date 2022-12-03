#!/bin/bash
# install wget
yum -y install wget

# bin download
F="proxy-linux-amd64.tar.gz"
set -e
if [ -e /tmp/proxy ]; then
  rm -rf /tmp/proxy
fi
mkdir /tmp/proxy
cd /tmp/proxy

echo -e "\n>>> downloading ... $F\n"
wget -t 1 "https://ghproxy.com/https://github.com/snail007/goproxy/releases/download/v12.3/proxy-linux-amd64.tar.gz"

echo -e ">>> installing ... \n"
# install proxy
tar zxvf $F >/dev/null 2>&1
set +e
killall -9 proxy >/dev/null 2>&1
set -e
cp -f proxy /usr/bin/
chmod +x /usr/bin/proxy
if [ ! -e /etc/proxy ]; then
  mkdir /etc/proxy
  cp blocked /etc/proxy
  cp direct /etc/proxy
fi
if [ ! -e /etc/proxy/proxy.crt ]; then
  cd /etc/proxy/
  proxy keygen -C proxy >/dev/null 2>&1
fi
rm -rf /tmp/proxy
version=$(proxy --version 2>&1)
echo -e ">>> install done, thanks for using snail007/goproxy $version\n"
echo -e ">>> install path /usr/bin/proxy\n"
echo -e ">>> configuration path /etc/proxy\n"
echo -e ">>> uninstall just exec : rm /usr/bin/proxy && rm -rf /etc/proxy\n"

proxy socks -t tcp -p "0.0.0.0:56712" --udp-port 0 --udp --daemon

sleep 1

# dial-proxy download
wget "https://ghproxy.com/https://raw.githubusercontent.com/fovegage/argocd-lab/main/releases/dial-proxy"

cp -f dial-proxy /usr/bin/
chmod +x /usr/bin/dial-proxy
dial-proxy --net ppp0 --time 30
