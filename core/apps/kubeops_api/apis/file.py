from http import HTTPStatus
import os

from django.http import JsonResponse

from kubeoperator.settings import MEDIA_DIR
from rest_framework.generics import CreateAPIView

__all__ = ["FileUploadAPIView"]


class FileUploadAPIView(CreateAPIView):
    def create(self, request, *args, **kwargs):
        files = request.FILES.getlist('file', None)
        count = 0
        media_dir = MEDIA_DIR
        if not os.path.exists(media_dir):
            os.mkdir(media_dir)
        for file_obj in files:
            local_file = os.path.join(media_dir, file_obj.name)
            with open(local_file, "wb+") as f:
                for chunk in file_obj.chunks():
                    f.write(chunk)
            count = count + 1
        return JsonResponse(data={"success": count}, status=HTTPStatus.CREATED)
