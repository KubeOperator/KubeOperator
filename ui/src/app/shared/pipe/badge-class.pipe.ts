import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'badgeClass'
})
export class BadgeClassPipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {
    let result = '';
    if (value) {
      switch (value) {
        case 'Warning':
          result = 'badge-danger';
          break;
        case 'WARNING':
          result = 'badge-danger';
          break;
        case 'INFO':
          result = '';
          break;
        case 'Normal':
          result = '';
          break;
      }
    }
    return result;
  }

}
