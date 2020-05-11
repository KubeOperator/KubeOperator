import {Pipe, PipeTransform} from '@angular/core';

@Pipe({name: 'keys'})
export class KeysPipe implements PipeTransform {
  transform(obj: Object, args: any[] = null): any {
    const ingore_keys = ['VC_USERNAME', 'VC_PASSWORD'];
    const array = [];
    Object.keys(obj).forEach(key => {
      if (ingore_keys.indexOf(key) === -1) {
        array.push({
          value: obj[key],
          key: key
        });
      }
    });
    return array;
  }
}
