import logging

from django.http import HttpResponse
from kombu.utils import json
from rest_framework.views import APIView

from ansible_api.permissions import IsSuperUser
from log.es import search_log

Logger = logging.getLogger(__name__)


class SearchSystemLog(APIView):
    permission_classes = (IsSuperUser,)

    def post(self, request, *args, **kwargs):
        logs = search_log(request.data)
        return HttpResponse(json.dumps(logs))
