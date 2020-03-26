import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'done'
})
export class DonePipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {

    let result = null;
    if (value) {
      switch (value) {
        case 'ENABLE':
          result = 'assets/images/done.svg';
          break;
        case 'DISABLE':
          result = 'assets/images/na.svg';
          break;
        default:
          result = '';
      }
    }
    return result;
  }

}
