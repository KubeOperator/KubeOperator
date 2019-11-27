import datetime

from elasticsearch import Elasticsearch
from elasticsearch_dsl import Search

from fit2ansible.settings import ELASTICSEARCH_HOST


def search_log(params):
    level = params.get('level', None)
    page = params.get('currentPage', None)
    keywords = params.get('keywords', None)
    limit_days = params.get('limit_days', None)

    size = 10
    time_start = get_start_time(limit_days)
    time_end = get_time_now()
    index = get_index()

    client = get_es_client()
    s = Search(using=client, index=index)
    s = s.using(client)
    if level and not level == 'all':
        s = s.query("match", levelname=level)
    s = s.query("range", timestamp={"gte": time_start, "lte": time_end})
    if page and size:
        s = s[(page - 1) * size:page * size]
    if keywords:
        s = s.query("match", message=keywords)
    s = s.sort({"timestamp": {"order": "desc"}})
    print(s.to_dict())
    s.execute()
    items = []
    for hit in s:
        items.append(
            {
                "name": hit.name,
                "level": hit.levelname,
                "timestamp": format_tz_time(hit.timestamp),
                "filename": hit.filename,
                "funcName": hit.funcName,
                "lineno": hit.lineno,
                "message": hit.message,
                "host_ip": hit.host_ip,
                "exc_text": hit.exc_text
            }
        )
    print(len(items))
    return {
        "items": items,
        "total": s.count()
    }


def format_tz_time(tz_time):
    _format = "%Y-%m-%dT%H:%M:%S.%fZ"
    format = "%Y-%m-%d %H:%M:%S"
    local_time = datetime.datetime.strptime(tz_time, _format) + datetime.timedelta(hours=8)
    return datetime.datetime.strftime(local_time, format)


def format_local_time(local_time):
    _format = "%Y-%m-%d %H:%M:%S"
    format = "%Y-%m-%dT%H:%M:%S.%fZ"
    tz_time = datetime.datetime.strptime(local_time, _format) - datetime.timedelta(hours=8)
    return datetime.datetime.strftime(tz_time, format)


def get_time_now():
    now = datetime.datetime.now()
    format = "%Y-%m-%dT%H:%M:%S.%fZ"
    return datetime.datetime.strftime(now, format)


def get_start_time(days):
    format = "%Y-%m-%dT%H:%M:%S.%fZ"
    time_start = datetime.datetime.now() - datetime.timedelta(days=int(days))
    return datetime.datetime.strftime(time_start, format)


def get_index():
    year = datetime.datetime.now().year
    month = datetime.datetime.now().month
    return 'kubeoperator-{}.{}'.format(year, month)


def get_es_client():
    client = Elasticsearch(hosts=[ELASTICSEARCH_HOST])
    return client


def index(client, index, doc_type, body):
    return client.index(index=index, doc_type=doc_type, body=body)


def update(client, index, doc_type, body, id):
    return client.update(index=index, doc_type=doc_type, body=body, id=id)


def exists(client, index):
    return client.indices.exists(index=index)
