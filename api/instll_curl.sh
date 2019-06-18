wget https://kubeoperator-public.oss-cn-beijing.aliyuncs.com/dependencies/curl-7.55.1.tar.gz -o /tmp/curl-7.55.1.tar.gz
cd /tmp
tar -xzvf  curl-7.55.1.tar.gz
cd curl-7.55.1
./configure
make
make install