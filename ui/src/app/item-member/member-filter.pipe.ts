import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'memberFilter',
  pure: false
})
export class MemberFilterPipe implements PipeTransform {

  transform(options: any[], values: any[]): any {
    options.forEach(o => {
      o['disabled'] = false;
      values.forEach(v => {
        if (o['value'] === v['value']) {
          console.log(o['text']);
          o['disable'] = true;
        }
      });
    });
    return options;
  }
}
