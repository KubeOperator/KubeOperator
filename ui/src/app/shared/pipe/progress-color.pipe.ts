import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'progressColor'
})
export class ProgressColorPipe implements PipeTransform {

  transform(value: any, ...args: any[]): any {
    let result = null;
    if (value) {
      if (value <= 60) {
        result = '#00af00';
        return result;
      } else if (value > 60 && value <= 80) {
        result = '#FFD700';
        return result;
      } else if (value > 80) {
        result = '#FF4040';
        return result;
      }
    }
    return result;
  }

}
