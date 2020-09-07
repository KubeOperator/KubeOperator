import {Pipe, PipeTransform} from '@angular/core';
import {TranslateService} from '@ngx-translate/core';

@Pipe({
    name: 'userType'
})
export class UserTypePipe implements PipeTransform {

    constructor(private translateService: TranslateService) {
    }

    transform(value: string, ...args: unknown[]): unknown {
        let result = '';
        if (value) {
            switch (value) {
                case 'LOCAL':
                    result = '' + this.translateService.instant('APP_USER_LOCAL');
                    break;
                case 'LDAP':
                    result = '' + this.translateService.instant('APP_USER_LDAP');
                    break;
                default:
                    result = value;
                    break;
            }
        }
        return result;
    }

}
