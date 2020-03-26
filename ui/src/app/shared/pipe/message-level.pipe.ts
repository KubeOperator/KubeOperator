import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'messageLevel'
})
export class MessageLevelPipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {
    let result = '';
    if (value) {
      switch (value) {
        case 'INFO':
          result = '信息';
          break;
        case 'WARNING':
          result = '告警';
          break;
        case 'ERROR':
          result = '错误';
          break;
      }
    }
    return result;
  }

}
