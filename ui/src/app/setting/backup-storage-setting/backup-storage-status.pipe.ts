import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'backupStorageStatus'
})
export class BackupStorageStatusPipe implements PipeTransform {

  transform(value: string): any {
    if (value === 'VALID') {
        return '可用';
    } else {
        return '不可用';
    }
  }

}
