FROM jpillora/dnsmasq:latest
ADD ./dnsmasq.conf /etc/dnsmasq.conf
ADD ./dnsmasq.d /etc/dnsmasq.d
EXPOSE 53/udp
EXPOSE 8080/tcp
ENV HTTP_USER=admin HTTP_PASS=DoNotChange
