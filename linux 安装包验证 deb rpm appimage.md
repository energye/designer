
---

deb: 展开所有文件
dpkg-deb -R packagename.deb ./packagename

---

rpm: 展开所有文件
rpm2cpio ./myapplcl_1.0.0.0_linux_amd64.rpm | cpio -idv

rpm: 查看包依赖
rpm -qpR myapplcl_1.0.0.0_linux_amd64.rpm

---

deb rpm appimage: 测试

先在宿主机开权限（必须执行）
xhost +local:root

deb: 测试
docker run -it --rm \
-e DISPLAY=$DISPLAY \
-v /tmp/.X11-unix:/tmp/.X11-unix \
-v $HOME/.Xauthority:/root/.Xauthority \
-v ${PWD}:/app \
amd64/ubuntu:22.04

更新: apt update 
安装deb: apt install xxx.deb 
运行: /usr/bin/xxx

appimage: 测试

rpm: 测试
docker run -it --rm \
-e DISPLAY=$DISPLAY \
-v /tmp/.X11-unix:/tmp/.X11-unix \
-v $HOME/.Xauthority:/root/.Xauthority \
-v ${PWD}:/app \
rockylinux/rockylinux:9

安装rpm: dnf install xxx.rpm
运行: /usr/bin/xxx
卸载: rpm -e --nodeps myapplcl
查看: rpm -qa | grep myapplcl