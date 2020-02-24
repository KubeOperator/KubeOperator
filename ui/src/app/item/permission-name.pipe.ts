import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'permissionName'
})
export class PermissionNamePipe implements PipeTransform {

  transform(value: string): any {
    if (value === 'VIEWER') {
      return '只读用户';
    } else if (value === 'MANAGER') {
      return '项目管理员';
    } else {
      return '超级管理员';
    }
  }
}
