# -*- coding: utf-8 -*-
#
import os
import copy
import tarfile
import zipfile
import gzip

import paramiko
from io import StringIO

from itsdangerous import TimedJSONWebSignatureSerializer, \
    JSONWebSignatureSerializer, BadSignature, SignatureExpired
from django.conf import settings


def get_object_or_none(model, **kwargs):
    try:
        obj = model.objects.get(**kwargs)
    except model.DoesNotExist:
        return None
    return obj


def shadow_key(iterable, key=None, remove=False):
    data_dump = copy.deepcopy(iterable)
    if isinstance(iterable, dict):
        for k, v in data_dump.items():
            if isinstance(key, (str, tuple)) and k == key or key(k):
                if remove:
                    del iterable[k]
                else:
                    iterable[k] = "******"
            if isinstance(v, (list, tuple, dict)):
                iterable[k] = shadow_key(v, key=key)
    elif isinstance(iterable, (list, tuple)):
        for item in data_dump:
            iterable.remove(item)
            iterable.append(shadow_key(item, key=key))
    return iterable


class Singleton(type):
    def __init__(cls, *args, **kwargs):
        cls.__instance = None
        super().__init__(*args, **kwargs)

    def __call__(cls, *args, **kwargs):
        if cls.__instance is None:
            cls.__instance = super().__call__(*args, **kwargs)
            return cls.__instance
        else:
            return cls.__instance


class Signer(metaclass=Singleton):
    """用来加密,解密,和基于时间戳的方式验证token"""
    def __init__(self, secret_key=None):
        self.secret_key = secret_key

    def sign(self, value):
        s = JSONWebSignatureSerializer(self.secret_key)
        return s.dumps(value).decode()

    def unsign(self, value):
        if value is None:
            return value
        s = JSONWebSignatureSerializer(self.secret_key)
        try:
            return s.loads(value)
        except BadSignature:
            return {}

    def sign_t(self, value, expires_in=3600):
        s = TimedJSONWebSignatureSerializer(self.secret_key, expires_in=expires_in)
        return str(s.dumps(value), encoding="utf8")

    def unsign_t(self, value):
        s = TimedJSONWebSignatureSerializer(self.secret_key)
        try:
            return s.loads(value)
        except (BadSignature, SignatureExpired):
            return {}


def get_signer():
    signer = Signer(settings.SECRET_KEY)
    return signer


def uncompress_tar(src_file, dest_dir):
    try:
        tar = tarfile.open(src_file)
        names = tar.getnames()
        for name in names:
            tar.extract(name, dest_dir)
        tar.close()
    except Exception as e:
        return False, e


def uncompress_zip(src_file, dest_dir):
    try:
        zip_file = zipfile.ZipFile(src_file)
        for names in zip_file.namelist():
            zip_file.extract(names, dest_dir)
        zip_file.close()
    except Exception as e:
        return False, e


def uncompress_gz(src_file, dest_dir):
    try:
        f_name = dest_dir + '/' + os.path.basename(src_file)
        # 获取文件的名称，去掉
        g_file = gzip.GzipFile(src_file)
        # 创建gzip对象
        open(f_name, "w+").write(g_file.read())
        g_file.close()
    except Exception as e:
        return False, e


def ssh_key_string_to_obj(text, password=None):
    key = None
    try:
        key = paramiko.RSAKey.from_private_key(StringIO(text), password=password)
    except paramiko.SSHException:
        pass

    try:
        key = paramiko.DSSKey.from_private_key(StringIO(text), password=password)
    except paramiko.SSHException:
        pass
    return key
