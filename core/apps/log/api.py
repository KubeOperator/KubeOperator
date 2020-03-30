import logging
from django.http import HttpResponse
from kombu.utils import json
from rest_framework.views import APIView
from log.log import SystemLog

Logger = logging.getLogger(__name__)

sys_log = SystemLog()


class SearchSystemLog(APIView):

    def post(self, request, *args, **kwargs):
        level = request.data.get('level', None)
        page = request.data.get('page', None)
        size = request.data.get('size', None)
        limit = request.data.get('limit', None)
        keywords = request.data.get('keywords', None)
        logs = sys_log.search(level, page, size, limit, keywords)
        return HttpResponse(json.dumps(logs))
