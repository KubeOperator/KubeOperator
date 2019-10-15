import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'clusterHealthStatus',
})
export class ClusterHealthStatusPipe implements PipeTransform {

  transform(value: string): any {
    if (value === '0') {
      return 'Running';
    } else {
      return 'Error';
    }
  }
}
