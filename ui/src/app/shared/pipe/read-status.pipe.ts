import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'readStatus'
})
export class ReadStatusPipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {
    let result = '';
    if (value) {
      switch (value) {
        case 'READ':
          result = '已读';
          break;
        case 'UNREAD':
          result = '未读';
          break;
      }
    }
    return result;
  }

}
