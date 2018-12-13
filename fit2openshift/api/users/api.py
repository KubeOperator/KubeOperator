# -*- coding: utf-8 -*-
#
from rest_framework.permissions import IsAuthenticated, IsAdminUser
from rest_framework.generics import RetrieveAPIView
from rest_framework.viewsets import ModelViewSet
from django.contrib.auth import get_user_model

from .serializers import ProfileSerializer, UserSerializer, UserCreateUpdateSerializer


class UserViewSet(ModelViewSet):
    queryset = get_user_model().objects.all()
    permission_classes = (IsAdminUser,)
    serializer_class = UserSerializer
    serializer_class_create = UserCreateUpdateSerializer

    def get_serializer_class(self):
        if self.action in ('create', 'update'):
            return self.serializer_class_create
        else:
            return super().get_serializer_class()


class UserProfileApi(RetrieveAPIView):
    permission_classes = (IsAuthenticated,)
    serializer_class = ProfileSerializer

    def get_object(self):
        return self.request.user
