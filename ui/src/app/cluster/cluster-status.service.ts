import {Injectable} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ClusterStatusService {

  constructor() {
  }

  // getLogo(status: string): string {
  //   let path = null;
  //   switch (status) {
  //     case 'UNKNOWN':
  //
  //   }
  // }

  getComment(status: string): string {
    let result = null;
    switch (status) {
      case 'UNKNOWN':
        result = '未知';
        break;
      case 'RUNNING':
        result = '运行中';
        break;
      case 'INSTALLING':
        result = '部署中';
        break;
      case 'ERROR':
        result = '错误';
        break;
      case 'WARNING':
        result = '警告';
        break;
    }
    console.log(result);
    return result;
  }
}
