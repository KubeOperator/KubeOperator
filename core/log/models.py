from django.db import models


# Create your models here.

class SystemLog():

    def __init__(self, name, timestamp, level, msg):
        self.name = name
        self.timestamp = timestamp
        self.level = level
        self.msg = msg



