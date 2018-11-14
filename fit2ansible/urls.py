from django.urls import include, path, re_path
from django.contrib import admin
from django.conf import settings
from django.conf.urls.static import static

from drf_yasg.views import get_schema_view
from drf_yasg import openapi
from . import error_handler


schema_view = get_schema_view(
   openapi.Info(
      title="Ansible UI Restful API",
      default_version='v1',
      description="It's ansible ui project restful api document",
      terms_of_service="http://www.jumpserver.org",
      contact=openapi.Contact(email="ibuler@fit2cloud.com"),
      license=openapi.License(name="GPLv2"),
   ),
   # validators=['flex', 'ssv'],
   public=True,
   # permission_classes=(permissions.AllowAny,),
)


def get_api_v1_urlpatterns():
    _urlpatterns = [
        path('', include('users.urls')),
        path('', include('ansible_api.urls.api_urls')),
        path('', include('celery_api.urls.api_urls')),
    ]
    return _urlpatterns


def get_view_url_patterns():
    from ansible_api.urls import view_urlpatterns as ansible_view_urlpatterns
    return ansible_view_urlpatterns


urlpatterns = [
    path('admin/', admin.site.urls),
    re_path(r'^docs(?P<format>\.json|\.yaml)/', schema_view.without_ui(cache_timeout=None), name='schema-json'),
    re_path(r'^swagger|docs/', schema_view.with_ui('swagger', cache_timeout=1), name='schema-swagger-ui'),
    path('redoc/', schema_view.with_ui('redoc', cache_timeout=None), name='schema-redoc'),
    path('api/v1/', include(get_api_v1_urlpatterns())),
    path('', include(get_view_url_patterns())),
]
urlpatterns += static(settings.STATIC_URL, document_root=settings.STATIC_ROOT)

handler404 = error_handler.error404
handler500 = 'rest_framework.exceptions.server_error'
handler400 = 'rest_framework.exceptions.bad_request'
