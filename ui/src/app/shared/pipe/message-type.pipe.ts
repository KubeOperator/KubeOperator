import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'messageType'
})
export class MessageTypePipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {
    let result = null;
    if (value) {
      switch (value) {
        case 'CLUSTER':
          result = '集群消息';
          break;
        case 'SYSTEM':
          result = '系统消息';
          break;
        default:
          result = '系统消息';
      }
    }
    return result;
  }

}
