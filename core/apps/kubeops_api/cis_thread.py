import threading


class CisThread(threading.Thread):

    def __init__(self, func):
        threading.Thread.__init__(self)
        self.func = func

    def run(self):
        self.func()
