from django.http import JsonResponse
from rest_framework import status


def error404(request, *args, **kwargs):
    data = {
        'error': 'Not found'
    }
    return JsonResponse(data, status=status.HTTP_400_BAD_REQUEST)
