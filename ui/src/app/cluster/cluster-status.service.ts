import {Injectable} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ClusterStatusService {

  constructor() {
  }


  getComment(status: string): string {
    let result = null;
    switch (status) {
      case 'READY':
        result = '准备安装';
        break;
      case 'RUNNING':
        result = '运行中';
        break;
      case 'INSTALLING':
        result = '部署中';
        break;
      case 'DELETING':
        result = '卸载中';
        break;
      case 'ERROR':
        result = '错误';
        break;
      case 'WARNING':
        result = '警告';
        break;
    }
    return result;
  }
}
