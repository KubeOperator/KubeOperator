import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'nullFilter'
})
export class NullFilterPipe implements PipeTransform {

  transform(value: any, args?: any): any {
    return value === null || value === '' ? '无' : value;
  }

}
