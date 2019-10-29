import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'status'
})
export class StatusPipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {
    let result = null;
    if (value) {
      switch (value) {
        case 'RUNNING':
          result = '运行中';
          break;
        case 'INSTALLING':
          result = '安装中';
          break;
        case 'UPGRADING':
          result = '升级中';
          break;
        case 'SCALING':
          result = '伸缩中';
          break;
        case 'ERROR':
          result = '错误';
          break;
        case 'DELETING':
          result = '卸载中';
          break;
        case 'READY':
          result = '就绪';
          break;
        case 'INITIALIZING':
          result = '初始化中';
          break;
        case 'UNKNOWN':
          result = '未知';
          break;
        case 'WARNING':
          result = '告警';
          break;
        case 'ADDING':
          result = '扩容中';
          break;
        default:
          result = '未知';
      }
    }
    return result;
  }

}
