import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'statusColor'
})
export class StatusColorPipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {
    let result = null;
    if (value) {
      switch (value) {
        case 'RUNNING':
          result = '#00af00';
          break;
        case 'WARNING':
          result = '#FFD700';
          break;
        case 'ERROR':
          result = '#FF4040';
          break;
        case 'UNKNOWN':
          result = '#575757';
          break;
        default:
          result = '#575757';
      }
    }
    return result;
  }

}
