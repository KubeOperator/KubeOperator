import jms_storage
import os


class StorageClient():

    def check_valid(self,backupStorage):
        storage_config = self.coverToConfig(self,backupStorage['credentials'])
        if 'S3' == backupStorage['type']:
            return self.check_s3(self,storage_config)


    def check_s3(self,storage_config):
        client= jms_storage.S3Storage(storage_config)
        # 上传文件测试可用性
        file_locale_path = os.path.abspath(os.path.join(os.getcwd(), os.path.pardir, 'README.md'))
        return client.is_valid(file_locale_path,'kube-operator-test')


    def coverToConfig(self,credentials):
        storage_config = {}
        storage_config['BUCKET'] = credentials.get('bucket',"kube-operator")
        storage_config['ACCESS_KEY'] = credentials.get('accessKey',None)
        storage_config['SECRET_KEY'] = credentials.get('secretKey',None)
        storage_config['CONTAINER_NAME'] = credentials.get('container',None)
        storage_config['ACCOUNT_NAME'] = credentials.get('accountName',None)
        storage_config['ACCOUNT_KEY'] = credentials.get('accountSecret',None)
        storage_config['ENDPOINT_SUFFIX'] = credentials.get('endpointSuffix',None)
        return storage_config