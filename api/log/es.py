import datetime

from elasticsearch import Elasticsearch
from elasticsearch_dsl import Search


# from fit2ansible.settings import ELASTICSEARCH_HOST


def search_log():
    level = "INFO"
    page = 1
    size = 10
    time_start = "2016-10-31T10:20:03"
    time_end = "2019-11-31T10:20:06"
    index = "my_python_app-2019.10.31"

    client = get_es_client()
    s = Search(using=client, index=index)
    s = s.using(client)
    s = s.query("match", levelname=level)
    s = s.query("range", timestamp={"gte": time_start, "lte": time_end})
    s = s[page:size]
    s = s.sort({"timestamp": {"order": "desc"}})
    s.execute()
    hits = []
    for hit in s:
        hits.append(
            {
                "name": hit.name,
                "level": hit.levelname,
                "timestamp": format_tz_time(hit.timestamp)
            }
        )
    return hits


def format_tz_time(tz_time):
    _format = "%Y-%m-%dT%H:%M:%S.%fZ"
    format = "%Y-%m-%d %H:%M:%S"
    local_time = datetime.datetime.strptime(tz_time, _format)
    return datetime.datetime.strftime(local_time, format)


def get_es_client():
    client = Elasticsearch(hosts=['172.16.10.142'])
    return client
