import datetime
import logging

from django.conf import settings
from elasticsearch import Elasticsearch, helpers
from elasticsearch_dsl import Search

logging.getLogger("requests").setLevel(logging.ERROR)
logging.getLogger("elasticsearch").setLevel(logging.ERROR)


def get_es_client():
    es_host = settings.ELASTICSEARCH_HOST
    return Elasticsearch(["{}.".format(es_host)])


def ensure_index_exists(client: Elasticsearch, name):
    if not client.indices.exists(name):
        client.indices.create(name)


def perform_query(search: Search):
    r = search.execute()
    return r.hits


def get_time_now():
    now = datetime.datetime.now()
    format = "%Y-%m-%dT%H:%M:%S.%fZ"
    return datetime.datetime.strftime(now, format)


def format_tz_time(tz_time):
    _format = "%Y-%m-%dT%H:%M:%S.%fZ"
    format = "%Y-%m-%d %H:%M:%S"
    local_time = datetime.datetime.strptime(tz_time, _format) + datetime.timedelta(hours=8)
    return datetime.datetime.strftime(local_time, format)


def get_start_time(days):
    format = "%Y-%m-%dT%H:%M:%S.%fZ"
    time_start = datetime.datetime.now() - datetime.timedelta(days=int(days))
    return datetime.datetime.strftime(time_start, format)


def get_index():
    date = datetime.datetime.now().strftime('%Y.%m')
    return 'kubeoperator-{}'.format(date)


def index(client, index, doc_type, body):
    return client.index(index=index, doc_type=doc_type, body=body)


def update(client, index, doc_type, body, id):
    return client.update(index=index, doc_type=doc_type, body=body, id=id)


def exists(client, index):
    return client.indices.exists(index=index)


def create_index_and_mapping(client, index, doc_type, mapping):
    client.indices.create(index=index)
    result = client.indices.put_mapping(doc_type=doc_type, index=index, body=mapping,
                                        params={'include_type_name': 'true'})
    return result['acknowledged']


def batch_data(client, data):
    return helpers.bulk(client, data)


def delete_index(client, index):
    return client.indices.delete(index=index)


def search_event(params, cluster_name):
    type = params.get('type', None)
    page = params.get('currentPage', None)
    size = params.get('size', None)
    keywords = params.get('keywords', None)
    limit_days = params.get('limitDays', None)
    time_start = get_start_time(limit_days)
    time_end = get_time_now()
    year = datetime.datetime.now().year
    month = datetime.datetime.now().month
    index = (cluster_name + '-{}.{}').format(year, month)

    client = get_es_client()
    s = Search(index=index).using(client)
    s = s.query("range", last_timestamp={"gte": time_start, "lte": time_end})
    if page and size:
        s = s[(page - 1) * size:page * size]
    if type and not type == 'all':
        s = s.query("match", type=type)
    if keywords:
        s = s.query("match", message=keywords)
    s = s.sort({"last_timestamp": {"order": "desc"}})
    s.execute()
    items = []
    for hit in s:
        items.append(
            {
                "uid": hit.uid,
                "action": hit.action,
                "type": hit.type,
                "last_timestamp": hit.last_timestamp,
                "cluster_name": hit.cluster_name,
                "component": hit.component,
                "host": hit.host,
                "message": hit.message,
                "first_timestamp": hit.first_timestamp,
                "name": hit.name,
                "namespace": hit.namespace,
                "reason": hit.reason,
            }
        )
    return {
        "items": items,
        "total": s.count()
    }


def get_event_uid_exist(client, index, uid):
    s = Search(index=index).using(client)
    s = s.query("match", uid=uid)
    s.execute()
    return s.count() == 0
