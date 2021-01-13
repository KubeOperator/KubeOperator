import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
    name: 'errFormat'
})
export class ErrFormatPipe implements PipeTransform {
    transform(value: string, ...args: unknown[]): unknown {
        if (value !== null) {
            let errItem = value;
            errItem = errItem.replace(/\\n/gi,'\n');
            errItem = errItem.replace(/\\u/gi,'%u');
            errItem = errItem.replace(/\\/gi,'');
            errItem = unescape(errItem)
            return errItem
        }
    }
}