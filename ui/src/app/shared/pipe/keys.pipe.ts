import {Pipe, PipeTransform} from '@angular/core';

@Pipe({name: 'keys'})
export class KeysPipe implements PipeTransform {
  transform(obj: Object, args: any[] = null): any {
    const array = [];
    Object.keys(obj).forEach(key => {
      array.push({
        value: obj[key],
        key: key
      });
    });
    return array;
  }
}
