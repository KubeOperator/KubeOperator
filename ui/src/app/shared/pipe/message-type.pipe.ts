import { Pipe, PipeTransform } from '@angular/core';
import {TranslateService} from '@ngx-translate/core';

@Pipe({
  name: 'messageType'
})
export class MessageTypePipe implements PipeTransform {

  constructor(private translateService: TranslateService) {

  }
    transform(value: any, ...args: any[]): any {
    let result = null;
    if (value) {
      switch (value) {
        case 'CLUSTER':
          result = this.translateService.instant('APP_MSG_CLUSTER');
          break;
        case 'SYSTEM':
          result = this.translateService.instant('APP_MSG_SYSTEM');
          break;
        default:
          result = this.translateService.instant('APP_MSG_SYSTEM');
      }
    }
    return result;
  }
}
