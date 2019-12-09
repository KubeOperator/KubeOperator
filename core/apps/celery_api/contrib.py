# -*- coding: utf-8 -*-
#
import sys
from celery.utils.log import LoggingProxy, _in_sighandler, safe_str


class NoStripLoggingProxy(LoggingProxy):
    def write(self, data):
        """Write message to logging object."""
        if _in_sighandler:
            return print(safe_str(data), file=sys.__stderr__)
        if getattr(self._thread, 'recurse_protection', False):
            # Logger is logging back to this file, so stop recursing.
            return
        data = data.rstrip()
        if data and not self.closed:
            self._thread.recurse_protection = True
            try:
                self.logger.log(self.loglevel, safe_str(data))
            finally:
                self._thread.recurse_protection = False
