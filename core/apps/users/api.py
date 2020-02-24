# -*- coding: utf-8 -*-
#
from django.contrib.auth.models import User
from django.http import HttpResponse
from kombu.utils import json
from rest_framework.generics import get_object_or_404, RetrieveUpdateAPIView
from rest_framework.permissions import IsAuthenticated
from rest_framework.views import APIView
from rest_framework.viewsets import ModelViewSet

from users.models import Profile
from .serializers import ProfileSerializer, UserSerializer, UserCreateUpdateSerializer


class UserViewSet(ModelViewSet):
    queryset = User.objects.all()
    serializer_class = UserSerializer
    serializer_class_create = UserCreateUpdateSerializer

    def get_serializer_class(self):
        if self.action in ('create', 'update'):
            return self.serializer_class_create
        else:
            return super().get_serializer_class()


class UserProfileRetrieveApi(RetrieveUpdateAPIView):
    permission_classes = (IsAuthenticated,)
    serializer_class = ProfileSerializer

    def get_object(self):
        obj = get_object_or_404(Profile, pk=self.request.user.profile.id)
        return obj


class UserPasswordChangeApi(APIView):
    def post(self, request, *args, **kwargs):
        pk = kwargs.get('pk')
        password = request.data.get('password')
        new_password = request.data.get('new_password')
        user = get_object_or_404(User, pk=pk)
        response = HttpResponse()
        response.write(json.dumps({'result': 'success'}))
        if user.check_password(password):
            if new_password:
                user.set_password(new_password)
                user.save()
            else:
                raise Exception('新密码不能为空！')
        else:
            raise Exception('原密码错误!')
        return response
