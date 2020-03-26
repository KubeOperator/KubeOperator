# -*- coding: utf-8 -*-
#

from django.conf.urls import url
from django.urls import path
from rest_framework.routers import DefaultRouter
from rest_framework_jwt.views import obtain_jwt_token, refresh_jwt_token

from . import api

app_name = "users"

router = DefaultRouter()
router.register(r'users', api.UserViewSet, 'user')
router.register(r'profiles', api.UserProfileViewSets, 'profiles')

urlpatterns = [
    url(r'^profile/$', api.UserProfileApi.as_view(), name='user-profile'),
    url(r'password/', api.UserPasswordChangeApi.as_view(), name='change-password'),
    url(r'token/auth/', obtain_jwt_token),
    url(r'token/refresh/', refresh_jwt_token),
    url(r'users/sync/', api.SyncUserFromLDAPApi.as_view(), name='sync-user')
]

urlpatterns += router.urls
