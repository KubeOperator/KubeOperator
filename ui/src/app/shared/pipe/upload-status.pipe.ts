import {Pipe, PipeTransform} from '@angular/core';
import {baseUrl} from '../../cluster/cluster-role.service';

@Pipe({
  name: 'uploadStatus'
})
export class UploadStatusPipe implements PipeTransform {

  transform(value: number): string {
    if (value !== null) {
      switch (value) {
        case 0:
          return '队列中';
        case 1:
          return '上传中';
        case 2:
          return '完成';
        case 3:
          return '取消';
      }
    }
    return null;
  }

}
