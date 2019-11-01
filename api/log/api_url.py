from django.conf.urls import url
from rest_framework.routers import DefaultRouter

from log import api

app_name = "log"

router = DefaultRouter()
urlpatterns = [
    url(r'^log/', api.SearchSystemLog.as_view(), name='log'),
]
