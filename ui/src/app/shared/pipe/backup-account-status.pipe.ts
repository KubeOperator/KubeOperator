import { Pipe, PipeTransform } from '@angular/core';
import {TranslateService} from '@ngx-translate/core';

@Pipe({
  name: 'backupAccountStatus'
})
export class BackupAccountStatusPipe implements PipeTransform {

  constructor(private translateService: TranslateService) {
  }

  transform(value: string, ...args: unknown[]): unknown {
    let result = '';
    if (value) {
      switch (value) {
        case 'VALID':
          result = '<img src="assets/images/done.svg" border-style="none" vertical-align="middle">'
              + this.translateService.instant('APP_STATUS_RUNNING');
          break;
        case 'INVALID':
          result = '<clr-icon style="color: red" shape="times"></clr-icon>' + this.translateService.instant('APP_STATUS_FAILED');
          break;
        default:
          result = value;
          break;
      }
    }
    return result;
  }
}
