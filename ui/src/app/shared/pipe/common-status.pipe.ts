import {Pipe, PipeTransform} from '@angular/core';
import {TranslateService} from '@ngx-translate/core';

@Pipe({
    name: 'commonStatus'
})
export class CommonStatusPipe implements PipeTransform {

    constructor(private translateService: TranslateService) {

    }

    transform(value: string, ...args: unknown[]): unknown {
        let result = '';
        if (value) {
            switch (value) {
                case 'Running':
                    result = '<img src="assets/images/done.svg" border-style="none" vertical-align="middle">'
                        + this.translateService.instant('APP_STATUS_RUNNING');
                    break;
                case 'Initializing':
                    result = this.translateService.instant('APP_STATUS_INITIALING');
                    break;
                case 'Synchronizing':
                    result = this.translateService.instant('APP_STATUS_SYNCHRONIZING');
                    break;
                case 'Creating':
                    result = this.translateService.instant('APP_STATUS_CREATING');
                    break;
                case 'NotReady':
                    result = this.translateService.instant('APP_STATUS_NOTREADY');
                    break;
                case 'Terminating':
                    result = this.translateService.instant('APP_STATUS_TERMINATING');
                    break;
                case 'Failed':
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
