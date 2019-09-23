import jms_storage
import os


class StorageClient():

    def check_valid(self,backupStorage):
        storage_config = self.coverToConfig(self,backupStorage['credentials'])
        client = {}
        if 'S3' == backupStorage['type']:
            client = jms_storage.S3Storage(storage_config)
        if 'OSS' == backupStorage['type']:
            client = jms_storage.OSSStorage(storage_config)
        if 'AZURE' == backupStorage['type']:
            client = jms_storage.AzureStorage(storage_config)
        if client is not None:
            # 上传文件测试可用性
            file_locale_path = os.path.abspath(os.path.join(os.getcwd(), os.path.pardir, 'README.md'))
            return client.is_valid(file_locale_path, 'kube-operator-test')


    def coverToConfig(self,credentials):
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