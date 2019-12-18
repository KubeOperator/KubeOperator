# -*- coding: utf-8 -*-
#

from rest_framework.generics import RetrieveAPIView
from django.conf import settings


class UserProfileApi(RetrieveAPIView):
    model = settings.AUTH_USER_MODEL
    serializer_class = None
