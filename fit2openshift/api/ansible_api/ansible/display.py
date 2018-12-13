# -*- coding: utf-8 -*-
#
from ansible.utils.color import stringc
from ansible.utils.display import Display


class LogFileDisplay(Display):
    def __init__(self, file_obj, **kwargs):
        self.file_obj = file_obj
        super().__init__(**kwargs)

    def _set_column_width(self):
        self.columns = 79

    def display(self, msg, color=None, stderr=False, screen_only=False, log_only=False):
        """ Display a message to the user

        Note: msg *must* be a unicode string to prevent UnicodeError tracebacks.
        """
        super().display(msg, color=color, stderr=stderr,
                        screen_only=screen_only, log_only=log_only)
        if log_only:
            return
        if color:
            msg = stringc(msg, color)
        if not msg.endswith(u'\n'):
            msg = msg + u'\n'
        self.file_obj.write(msg)
        self.file_obj.flush()
