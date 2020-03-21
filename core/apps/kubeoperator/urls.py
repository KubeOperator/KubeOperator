from django.urls import include, path, re_path
from django.contrib import admin
from django.conf.urls.static import static
from drf_yasg.views import get_schema_view
from drf_yasg import openapi

from kubeoperator import settings
from kubeoperator.celery_flower import celery_flower_view
from . import error_handler

schema_view = get_schema_view(
    openapi.Info(
        title="KubeOperator Restful API",
        default_version='v1',
        terms_of_service="http://www.kubeoperator.io",
        contact=openapi.Contact(email="support@fit2cloud.com"),
        license=openapi.License(name="Apache 2.0"),
    ),
    public=True,
)


def get_api_v1_urlpatterns():
    _urlpatterns = [
        path('', include('users.urls')),
        path('', include('celery_api.urls.api_urls')),
        path('', include('kubeops_api.api_url')),
        path('', include('cloud_provider.api_url')),
        path('', include('storage.api_url')),
        path('', include('log.api_url')),
        path('', include('message_center.api_url')),

    ]
    return _urlpatterns


def get_view_url_patterns():
    from ansible_api.urls import view_urlpatterns as ansible_view_urlpatterns
    return ansible_view_urlpatterns


urlpatterns = [
    path('admin/', admin.site.urls),
    re_path(r'^docs(?P<format>\.json|\.yaml)/', schema_view.without_ui(cache_timeout=None), name='schema-json'),
    re_path(r'^swagger|docs/', schema_view.with_ui('swagger', cache_timeout=1), name='schema-swagger-ui'),
    re_path(r'flower/(?P<path>.*)', celery_flower_view, name='flower-view'),
    path('redoc/', schema_view.with_ui('redoc', cache_timeout=None), name='schema-redoc'),
    path('api/v1/', include(get_api_v1_urlpatterns())),
    path('', include(get_view_url_patterns())),
]
urlpatterns += static(settings.STATIC_URL, document_root=settings.STATIC_ROOT)

handler404 = error_handler.error404
handler500 = 'rest_framework.exceptions.server_error'
handler400 = 'rest_framework.exceptions.bad_request'
