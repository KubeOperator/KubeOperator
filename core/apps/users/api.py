# -*- coding: utf-8 -*-
#
from django.contrib.auth.models import User
from rest_framework import status, viewsets
from rest_framework.generics import get_object_or_404, UpdateAPIView, RetrieveAPIView, CreateAPIView
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.viewsets import ModelViewSet

from users.models import Profile
from .serializers import ProfileSerializer, UserSerializer, UserCreateUpdateSerializer, ChangeUserPasswordSerializer
from message_center.models import UserNotificationConfig, UserReceiver, UserMessage
from .tasks import start_sync_user_form_ldap
from rest_framework.exceptions import ValidationError


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

        config = UserNotificationConfig()
        config.create_config_by_username(username=serializer.data['username'])

        return Response(serializer.data, status=status.HTTP_201_CREATED, headers=headers)


class UserProfileApi(RetrieveAPIView):
    permission_classes = (IsAuthenticated,)
    serializer_class = ProfileSerializer

    def get_object(self):
        obj = get_object_or_404(Profile, pk=self.request.user.profile.id)
        return obj


class UserProfileViewSets(viewsets.ModelViewSet):
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
        try:
            instance = super().update(request, *args, **kwargs)
            if instance:
                return Response({"result": "ok"}, status=status.HTTP_202_ACCEPTED)
        except ValidationError as e:
            return Response(data={'msg': '原密码错误！'}, status=status.HTTP_400_BAD_REQUEST)


class SyncUserFromLDAPApi(CreateAPIView):
    # permission_classes = (IsSuperUser,)
    def post(self, request, *args, **kwargs):
        start_sync_user_form_ldap.apply_async()
        return Response({"result": "ok"}, status=status.HTTP_201_CREATED)
