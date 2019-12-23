from django.db import models

__all__ = ['Condition','HealthChecker','HealthCheck']


class Condition(models.Model):
    type = models.CharField(max_length=128, null=True)
    status = models.BooleanField(default=True, null=True)
    message = models.CharField(max_length=256, null=True)
    reason = models.CharField(max_length=256, null=True)
    last_time = models.DateTimeField(auto_now_add=True)


class HealthCheck(models.Model):
    def run(self):
        pass


class HealthChecker():
    def check(self):
        pass
