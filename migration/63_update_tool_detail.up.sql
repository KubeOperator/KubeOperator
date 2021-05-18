DELETE FROM ko_cluster_tool_detail WHERE name IN ('grafana', 'loki', 'prometheus', 'dashboard', 'registry', 'chartmuseum') AND architecture='amd64';
UPDATE ko_cluster_tool_detail SET architecture='all' WHERE name IN ('grafana', 'loki', 'prometheus', 'dashboard', 'registry', 'chartmuseum') AND architecture='arm64';

UPDATE ko_cluster_tool_detail SET vars='{\"grafana_image_name\":\"grafana/grafana\",\"grafana_image_tag\":\"7.3.3\",\"busybox_image_name\":\"busybox\",\"busybox_image_tag\":\"1.28\",\"curl_image_name\":\"curlimages/curl\",\"curl_image_tag\":\"7.73.0\"}' WHERE name='grafana' AND version='v7.3.3';
UPDATE ko_cluster_tool_detail SET vars='{\"loki_image_name\":\"grafana/loki\",\"loki_image_tag\":\"2.0.0\",\"promtail_image_name\":\"grafana/promtail\",\"promtail_image_tag\":\"2.0.0\"}' WHERE name='loki' AND version='v2.0.0';
UPDATE ko_cluster_tool_detail SET vars='{\"chartmuseum_image_name\":\"kubeoperator/chartmuseum\",\"chartmuseum_image_tag\":\"v0.12.0\"}' WHERE name='chartmuseum' AND version='v0.12.0';
UPDATE ko_cluster_tool_detail SET vars='{\"registry_image_name\":\"registry\",\"registry_image_tag\":\"2.7.1\"}' WHERE name='registry' AND version='v2.7.1';