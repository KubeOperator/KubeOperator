import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'itemRole'
})
export class ItemRolePipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {
    if (value) {
      if (value === 'VIEWER') {
        return '只读用户';
      }
      if (value === 'MANAGER') {
        return '管理员';
      }
    }
    return '未知';
  }

}
