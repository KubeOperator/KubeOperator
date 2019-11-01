from django.http import HttpResponse
from kombu.utils import json
from rest_framework.views import APIView

from ansible_api.permissions import IsSuperUser
from log.es import search_log


class SearchSystemLog(APIView):
    permission_classes = (IsSuperUser,)

    def post(self, request, *args, **kwargs):
        logs = search_log()
        return HttpResponse(json.dumps(logs))
