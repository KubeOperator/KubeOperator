import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'messageDetail'
})
export class MessageDetailPipe implements PipeTransform {

  transform(value: string, ...args: any[]): any {
    let result = null;
    if (value) {
        result = value.replace('>', '<br>');
    }
    return result;
  }

}
