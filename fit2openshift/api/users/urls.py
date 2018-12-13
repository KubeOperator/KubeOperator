# -*- coding: utf-8 -*-
#

from django.conf.urls import url
from rest_framework.routers import DefaultRouter
from . import api

app_name = "users"

router = DefaultRouter()
router.register(r'users', api.UserViewSet, 'user')


urlpatterns = [
    url(r'^profile/$', api.UserProfileApi.as_view(), name='user-profile'),
]

urlpatterns += router.urls
