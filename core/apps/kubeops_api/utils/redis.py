import redis
from django.conf import settings


class RedisHelper:
    __conn = redis.StrictRedis(host=settings.REDIS_HOST, port=settings.REDIS_PORT)

    def publish(self, channel, msg):
        self.__conn.publish(channel, msg)

    def subscribe(self, channel):
        pub = self.__conn.pubsub()
        pub.subscribe(channel)
        pub.parse_response()
        return pub
