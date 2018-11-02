import time
import threading
from channels.generic.websocket import JsonWebsocketConsumer


class CeleryLogWebsocket(JsonWebsocketConsumer):
    task = ''
    task_log_f = None
    disconnected = False

    def connect(self):
        task_id = self.scope['url_route']['kwargs']['task_id']
        try:
            self.task = CeleryTask.objects.get(id=task_id)
        except CeleryTask.DoesNotExist:
            self.send({'message': "Task {} not found".format(task_id)})
            self.disconnect(None)
            return
        try:
            self.task_log_f = open(self.task.log_path)
        except OSError:
            self.send({'message': "Task {} log not found".format(task_id)})
            self.disconnect(None)
            return

        self.accept()
        self.send_log_to_client()

    def disconnect(self, close_code):
        self.disconnected = True
        if self.task_log_f and not self.task_log_f.closed:
            self.task_log_f.close()
        self.close()

    def send_log_to_client(self):
        def func():
            while not self.disconnected:
                data = self.task_log_f.read(4096)
                if data:
                    self.send_json({'message': data})
                time.sleep(0.2)
        thread = threading.Thread(target=func)
        thread.start()
