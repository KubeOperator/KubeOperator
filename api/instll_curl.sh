wget https://curl.haxx.se/download/curl-7.55.1.tar.gz -o /tmp
cd /tmp
tar -xzvf  curl-7.55.1.tar.gz
cd curl-7.55.1
./configure
make
make install