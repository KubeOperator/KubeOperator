import jms_storage
import os
from kubeops_api.models.backup_storage import BackupStorage;


class StorageClient():

    def __init__(self,backupStorage):
        if isinstance(backupStorage,BackupStorage):
            backupStorage=backupStorage.get_dict()
        storage_config = self.cover_to_config(backupStorage['credentials'])
        try:
            if 'S3' == backupStorage['type']:
                self.client = jms_storage.S3Storage(storage_config)
            if 'OSS' == backupStorage['type']:
                self.client = jms_storage.OSSStorage(storage_config)
            if 'AZURE' == backupStorage['type']:
                self.client = jms_storage.AzureStorage(storage_config)
        except ValueError:
            pass

    def check_valid(self):
        if self.client is None:
            return False
        # 上传文件测试可用性
        return self.client.is_valid("../Dockerfile", 'kube-operator-test')

    def cover_to_config(self,credentials):
        storage_config = {}
        storage_config['BUCKET'] = credentials.get('bucket',"kube-operator")
        storage_config['ACCESS_KEY'] = credentials.get('accessKey',None)
        storage_config['SECRET_KEY'] = credentials.get('secretKey',None)
        storage_config['CONTAINER_NAME'] = credentials.get('container',None)
        storage_config['ACCOUNT_NAME'] = credentials.get('accountName',None)
        storage_config['ACCOUNT_KEY'] = credentials.get('accountSecret',None)
        storage_config['ENDPOINT_SUFFIX'] = credentials.get('endpointSuffix',None)
        storage_config['ENDPOINT'] = credentials.get('endpoint',None)
        return storage_config

    def list_buckets(self):
        return self.client.list_buckets()

    def upload_file(self,src,target):
        return self.client.upload(src,target)

    def exists(self,path):
        return self.client.exists(path)

    def download_file(self,src,target):
        return self.client.download(src,target)



