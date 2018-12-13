import logging


def get_logger(name):
    return logging.getLogger('ansible_api.%s' % name)
