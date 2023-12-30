FROM envoyproxy/envoy:v1.18-latest@sha256:0f1c9e46ce1d07faea29e768e703aef50b43564434ed0eb4a195c8d4e8a08327

COPY envoy.yaml /etc/envoy/envoy.yaml

CMD /usr/local/bin/envoy -c /etc/envoy/envoy.yaml -l trace --log-path /tmp/envoy_info.log
