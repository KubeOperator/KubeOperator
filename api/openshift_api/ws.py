import threading
import time
from channels.generic.websocket import JsonWebsocketConsumer

from openshift_api.models.deploy import DeployExecution


class F2OWebsocket(JsonWebsocketConsumer):
    disconnected = False
    execution_id = None

    def connect(self):
        self.execution_id = self.scope['url_route']['kwargs']['execution_id']
        if self.execution_id is None:
            raise Exception('execution_id not be None!')
        self.accept()
        self.send_deploy_execution()

    def send_deploy_execution(self):
        def func():
            last_data = DeployExecution.objects.filter(id=self.execution_id).first().to_json()
            while not self.disconnected:
                data = DeployExecution.objects.filter(id=self.execution_id).first().to_json()
                if data and not data == last_data:
                    self.send_json({'message': data})
                    last_data = data
                time.sleep(2)

        thread = threading.Thread(target=func)
        thread.start()

    def disconnect(self, close_code):
        self.disconnected = True
        self.close()
