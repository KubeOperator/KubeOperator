import {Pipe, PipeTransform} from '@angular/core';
import {TranslateService} from '@ngx-translate/core';

@Pipe({
    name: 'emailShow'
})
export class EmailShowPipe implements PipeTransform {

    constructor(private translateService: TranslateService) {

    }

    transform(value: string, ...args: unknown[]): unknown {
        let result = '';
        if (value.indexOf('@') === -1 || value.indexOf('.') === -1) {
            return result
        }
        const aiteIndex = value.indexOf('@')
        const pointIndex = value.lastIndexOf('.')
        const mail = value.substring(0, aiteIndex)
        if (mail.length <= 3) {
            result += '***'
        } else {
            result += value.substring(0, 3) + '***'
        }
        result += value.substring(pointIndex+1, value.length)
        return result;
    }
}
