"""
KubeOperator  配置文件
"""

import os
import datetime
from celery.schedules import crontab
from cmreslogging.handlers import CMRESHandler
from .conf import load_user_config

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
ANSIBLE_PROJECTS_DIR = os.path.join(BASE_DIR, 'data', 'ansible', 'projects')
TERRAFORM_DIR = os.path.join(BASE_DIR, 'data', 'terraform', 'projects')
BASE_LOG_DIR = os.path.join(BASE_DIR, "data", "log")
MEDIA_DIR = os.path.join(BASE_DIR, "data", "media")
VERSION_DIR = os.path.join(BASE_DIR, "build", "version")
CLOUDS_RESOURCE_DIR = os.path.join(BASE_DIR, "resource", "clouds")
CLUSTER_CONFIG_DIR = os.path.join(BASE_DIR, "resource", "cluster")
KUBEEASZ_DIR = os.path.join(BASE_DIR, "resource", "kubeasz")
WEBKUBECTL_URL = "http://webkubectl:8080/api/kube-config"
PACKAGE_IMAGE_NAME = 'registry.fit2cloud.com/public/nexus-helm:3.15.2-01'
PACKAGE_PATH_PREFIX = "/opt/kubeoperator/data/packages/"
PACKAGE_DIR = "/data/packages"
CONFIG = load_user_config()

SECRET_KEY = CONFIG.SECRET_KEY
DEBUG = CONFIG.DEBUG
ALLOWED_HOSTS = ['*']

INSTALLED_APPS = [
    'message_center.apps.MessageCenterConfig',
    'storage.apps.StorageConfig',
    'kubeops_api.apps.KubeOperatorApiConfig',
    'cloud_provider.apps.CloudProviderConfig',
    'ansible_api.apps.AnsibleApiConfig',
    'celery_api.apps.CeleryApiConfig',
    'log.apps.LogConfig',
    'users.apps.UsersConfig',
    'django_celery_beat',
    'rest_framework',
    'drf_yasg',
    'channels',
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
]

MIDDLEWARE = [
    'django.middleware.security.SecurityMiddleware',
    'django.contrib.sessions.middleware.SessionMiddleware',
    'django.middleware.common.CommonMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'django.middleware.clickjacking.XFrameOptionsMiddleware',
]

ROOT_URLCONF = 'kubeoperator.urls'

WSGI_APPLICATION = 'kubeoperator.wsgi.application'
ASGI_APPLICATION = 'kubeoperator.routing.application'

TEMPLATES = [
    {
        'BACKEND': 'django.template.backends.django.DjangoTemplates',
        'DIRS': [os.path.join(os.path.dirname(__file__), 'templates').replace('\\', '/'), ],
        'APP_DIRS': True,
        'OPTIONS': {
            'context_processors': [
                'django.template.context_processors.debug',
                'django.template.context_processors.request',
                'django.contrib.auth.context_processors.auth',
                'django.contrib.messages.context_processors.messages',
            ],
        },
    },
]

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.mysql',
        'NAME': CONFIG.DB_NAME,
        'USER': CONFIG.DB_USER,
        'PASSWORD': CONFIG.DB_PASSWORD,
        'HOST': CONFIG.DB_HOST,
        'PORT': CONFIG.DB_PORT
    }
}

AUTH_PASSWORD_VALIDATORS = [
    {
        'NAME': 'django.contrib.auth.password_validation.UserAttributeSimilarityValidator',
    },
    {
        'NAME': 'django.contrib.auth.password_validation.MinimumLengthValidator',
    },
    {
        'NAME': 'django.contrib.auth.password_validation.CommonPasswordValidator',
    },
    {
        'NAME': 'django.contrib.auth.password_validation.NumericPasswordValidator',
    },
]

STATIC_URL = '/static/'
STATIC_ROOT = os.path.join(BASE_DIR, "data", "static")

REDIS_HOST = CONFIG.REDIS_HOST
REDIS_PORT = CONFIG.REDIS_PORT
REDIS_PASSWORD = CONFIG.REDIS_PASSWORD

LOGIN_URL = '/admin/login'

CACHES = {
    'default': {
        'BACKEND': 'redis_cache.RedisCache',
        'LOCATION': 'redis://:%(password)s@%(host)s:%(port)s/4' % {
            'password': REDIS_PASSWORD,
            'host': REDIS_HOST,
            'port': REDIS_PORT,
        }
    },
}

CELERY_LOG_BROKER_URL = 'redis://:%(password)s@%(host)s:%(port)s/10' % {
    'password': REDIS_PASSWORD,
    'host': REDIS_HOST,
    'port': REDIS_PORT,
}

# Dump all celery log to here
CELERY_LOG_DIR = os.path.join(BASE_DIR, 'data', 'celery')

# Celery using redis as broker
CELERY_BROKER_URL = 'redis://:%(password)s@%(host)s:%(port)s/11' % {
    'password': REDIS_PASSWORD,
    'host': REDIS_HOST,
    'port': REDIS_PORT,
}
CELERY_TASK_SERIALIZER = 'pickle'
CELERY_RESULT_SERIALIZER = 'pickle'
CELERY_RESULT_BACKEND = CELERY_BROKER_URL
CELERY_ACCEPT_CONTENT = ['json', 'pickle']
CELERY_RESULT_EXPIRES = 3600
CELERY_WORKER_TASK_LOG_FORMAT = '%(message)s'
CELERY_WORKER_LOG_FORMAT = '%(message)s'
CELERY_TASK_EAGER_PROPAGATES = True
CELERY_WORKER_REDIRECT_STDOUTS = True
CELERY_WORKER_REDIRECT_STDOUTS_LEVEL = "ERROR"
FLOWER_URL = "localhost:5555"

REST_FRAMEWORK = {
    # Use Django's standard `django.contrib.auth` permissions,
    # or allow read-only access for unauthenticated users.
    'DEFAULT_FILTER_BACKENDS': (
        'django_filters.rest_framework.DjangoFilterBackend',
        'rest_framework.filters.SearchFilter',
        'rest_framework.filters.OrderingFilter',
    ),
    'DEFAULT_AUTHENTICATION_CLASSES': (
        'rest_framework_jwt.authentication.JSONWebTokenAuthentication',
        'rest_framework.authentication.SessionAuthentication',
    ),
    'ORDERING_PARAM': "order",
    'SEARCH_PARAM': "search",
    # 'DEFAULT_PAGINATION_CLASS': 'rest_framework.pagination.LimitOffsetPagination',
    # 'PAGE_SIZE': 25
}

SWAGGER_SETTINGS = {
    'DEFAULT_AUTO_SCHEMA_CLASS': 'kubeoperator.swagger.CustomSwaggerAutoSchema',
}

CHANNEL_REDIS = "redis://:{}@{}:{}/0".format(
    CONFIG.REDIS_PASSWORD, CONFIG.REDIS_HOST, CONFIG.REDIS_PORT
)
CHANNEL_LAYERS = {
    'default': {
        'BACKEND': 'channels_redis.core.RedisChannelLayer',
        'CONFIG': {
            "hosts": [CHANNEL_REDIS],
        },
    },
}

JWT_AUTH = {
    'JWT_EXPIRATION_DELTA': datetime.timedelta(hours=12),
    'JWT_AUTH_HEADER_PREFIX': 'JWT',
    'JWT_RESPONSE_PAYLOAD_HANDLER': 'users.utils.jwt_response_payload_handler',
    'JWT_ALLOW_REFRESH': True,
}

ELASTICSEARCH_HOST = CONFIG.ELASTICSEARCH_HOST
ELASTICSEARCH_PORT = CONFIG.ELASTICSEARCH_PORT

LOGGING = {
    'version': 1,
    'disable_existing_loggers': False,
    'formatters': {
        'default': {
            'format': '[%(levelname)s][%(asctime)s][%(filename)s:%(lineno)d]%(message)s'
        },
    },
    'filters': {
        'require_debug_true': {
            '()': 'django.utils.log.RequireDebugTrue',
        },
    },
    'handlers': {
        'console': {
            'level': 'DEBUG',
            'filters': ['require_debug_true'],
            'class': 'logging.StreamHandler',
            'formatter': 'default'
        },
        'elasticsearch': {
            'level': 'INFO',
            'class': 'cmreslogging.handlers.CMRESHandler',
            'hosts': [{'host': ELASTICSEARCH_HOST, 'port': ELASTICSEARCH_PORT}],
            'es_index_name': 'ko-log',
            'index_name_frequency': CMRESHandler.IndexNameFrequency.MONTHLY,
            'auth_type': CMRESHandler.AuthType.NO_AUTH,
            'use_ssl': False,
        },
    },
    'loggers': {
        "": {
            'handlers': ['console'],
            'level': 'INFO',
        },
        'user': {
            'handlers': ['console', 'elasticsearch'],
            'level': 'INFO',
        },
        'cluster': {
            'handlers': ['console', 'elasticsearch'],
            'level': 'INFO',
        },
        'host': {
            'handlers': ['console', 'elasticsearch'],
            'level': 'INFO',
        },
        'log': {
            'handlers': ['console', 'elasticsearch'],
            'level': 'INFO',
        },
        'cloud_provider': {
            'handlers': ['console', 'elasticsearch'],
            'level': 'INFO',
        },
    },
}
AUTHENTICATION_BACKENDS = [
    'django.contrib.auth.backends.ModelBackend', 'users.authentication.ldap.LDAPAuthorizationBackend'
]

AUTH_LDAP_ENABLE = False
AUTH_LDAP_START_TLS = False
AUTH_LDAP_CONNECT_TIMEOUT = 30,
AUTH_LDAP_SEARCH_PAGED_SIZE = 10000,
AUTH_LDAP_SYNC_IS_PERIODIC = False
AUTH_LDAP_SYNC_INTERVAL = None
AUTH_LDAP_SYNC_CRONTAB = None
AUTH_LDAP_USER_LOGIN_ONLY_IN_USERS = False
AUTH_LDAP_OPTIONS_OPT_REFERRALS = -1

NODE_CREDENTIAL = {
    'username': "root",
    'password': "KubeOperator@2019"
}
