import logging
from datetime import datetime, timedelta
from elasticsearch_dsl import Search

from log.es import get_es_client, ensure_index_exists, perform_query, format_tz_time

logger = logging.getLogger(__name__)

levels = {
    "DEBUG": 3,
    "INFO": 2,
    "WARNING": 1,
    "ERROR": 0
}


class SystemLog:
    def __init__(self):
        self.index_name = 'ko-log-{}'.format(datetime.now().strftime('%Y.%m'))
        self.client = get_es_client()
        ensure_index_exists(self.client, self.index_name)

    def search(self, level, page, size, limit, keywords=None):
        s = self.build_query(level, page, size, limit, keywords)
        hits = perform_query(s)
        items = parse_data(hits)
        count = s.count()
        result = {"items": items, "total": count}
        return result

    def build_query(self, level, page, size, limit, keywords=None):
        s = Search(using=self.client, index=self.index_name)
        if level:
            ls = []
            for k in levels:
                if levels[k] <= levels[level]:
                    ls.append(k.lower())
            s = s.query("terms", levelname=ls)
        if page and size:
            s = s[(page - 1) * size:page * size]
        if keywords:
            s = s.query("match", message=keywords)
        if limit:
            now = datetime.now()
            start_time = now - timedelta(days=int(limit))
            s = s.query("range", timestamp={"gte": format_date(start_time), "lte": format_date(now)})
        return s


def format_date(date):
    formatter = "%Y-%m-%d"
    return datetime.strftime(date, formatter)


def parse_data(hits):
    items = []
    for hit in hits:
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
    return items
