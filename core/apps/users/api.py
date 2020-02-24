# -*- coding: utf-8 -*-
#
from django.contrib.auth.models import User
from rest_framework.generics import get_object_or_404, RetrieveUpdateAPIView, UpdateAPIView
from rest_framework.permissions import IsAuthenticated
from rest_framework.viewsets import ModelViewSet
from rest_framework.response import Response

from ansible_api.permissions import IsSuperUser
from users.models import Profile
from .serializers import ProfileSerializer, UserSerializer, UserCreateUpdateSerializer, ChangeUserPasswordSerializer
from rest_framework import status, viewsets


class UserViewSet(ModelViewSet):
    queryset = User.objects.all()
    serializer_class = UserSerializer
    serializer_class_create = UserCreateUpdateSerializer

    def get_serializer_class(self):
        if self.action in ('create', 'update'):
            return self.serializer_class_create
        else:
            return super().get_serializer_class()

    def create(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        self.perform_create(serializer)
        headers = self.get_success_headers(serializer.data)
        serializer.data.pop("password")
        return Response(serializer.data, status=status.HTTP_201_CREATED, headers=headers)


class UserProfileApi(RetrieveUpdateAPIView):
    permission_classes = (IsAuthenticated,)
    serializer_class = ProfileSerializer

    def get_object(self):
        obj = get_object_or_404(Profile, pk=self.request.user.profile.id)
        return obj


class UserProfileViewSets(viewsets.ModelViewSet):
    permission_classes = (IsSuperUser,)
    serializer_class = ProfileSerializer
    queryset = Profile.objects.all()
    lookup_url_kwarg = "id"
    lookup_field = "id"


class UserPasswordChangeApi(UpdateAPIView):
    permission_classes = (IsAuthenticated,)
    serializer_class = ChangeUserPasswordSerializer

    def get_object(self):
        obj = get_object_or_404(User, pk=self.request.user.id)
        return obj

    http_method_names = ["put", "option", "head"]

    def update(self, request, *args, **kwargs):
        instance = super().update(request, *args, **kwargs)
        if instance:
            return Response({"result": "ok"}, status=status.HTTP_202_ACCEPTED)
