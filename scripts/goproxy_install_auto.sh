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

# start socks proxy
proxy socks -t tcp -p "0.0.0.0:56712" --udp-port 0 --udp --daemon

# dial-proxy download
cd /root/
dialName="dial-proxy"
if [ -f "$dialName" ]; then
  rm -rf "$dialName"
fi

# kill -9 `cat dial-proxy.pid`
#dialPid="dial-proxy.pid"
#if [ -f "$dialPid" ]; then
#  # shellcheck disable=SC2046
#  # shellcheck disable=SC2006
#  kill -9 $(cat "$dialPid")
#  rm -rf "$dialPid"
#  # shellcheck disable=SC2105
#  continue
#fi

wget "https://ghproxy.com/https://raw.githubusercontent.com/fovegage/releases/main/dial-proxy"
cp -f /root/dial-proxy /usr/bin/
chmod +x /usr/bin/dial-proxy

# start dial proxy
# shellcheck disable=SC2028
nohup dial-proxy --net ppp --time 30 >nohup.out &
