import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'subscribeCheck'
})
export class SubscribeCheckPipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {
    let result = null;
    if (value) {
      switch (value) {
        case 'ENABLE':
          result = true;
          break;
        case 'DISABLE':
          result = false;
          break;
        default:
          result = false;
      }
    }
    return result;
  }

}
