# -*- coding: utf-8 -*-
#

"""
Immediately run adhoc and playbook
"""

from rest_framework import permissions, generics
from rest_framework.response import Response

from ..serializers import IMPlaybookSerializer, IMAdHocSerializer
<<<<<<< HEAD
from ..tasks import execute_playbook, run_im_adhoc
=======
from ..tasks import execute_playbook, run_adhoc_raw
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce

__all__ = ['IMPlaybookApi', 'IMAdHocApi']


class IMPlaybookApi(generics.CreateAPIView):
    permission_classes = (permissions.AllowAny,)
    serializer_class = IMPlaybookSerializer

    def create(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        if serializer.is_valid():
            playbook = serializer.save()
            task = execute_playbook.delay(playbook.id, save_history=False)
            return Response({'task': task.id})
        else:
            return Response({"error": serializer.errors})


class IMAdHocApi(generics.CreateAPIView):
    permission_classes = (permissions.AllowAny,)
    serializer_class = IMAdHocSerializer

    def create(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        if serializer.is_valid():
<<<<<<< HEAD
            task = run_im_adhoc.delay(
=======
            task = run_adhoc_raw.delay(
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
                serializer.validated_data.get('adhoc'),
                serializer.validated_data.get('inventory'),
            )
            return Response({'task': task.id})
        else:
            return Response({"error": serializer.errors})
