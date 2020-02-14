import {Pipe, PipeTransform} from '@angular/core';
import {ItemResourceDTO} from './item-resource';

@Pipe({
  name: 'itemResource'
})
export class ItemResourcePipe implements PipeTransform {

  transform(itemResources: ItemResourceDTO[], resourceType: string): any {
    if (itemResources !== undefined) {
      const choose = [];
      for (const i of itemResources) {
        if (i['resource_type'] === resourceType) {
          choose.push(i);
        }
      }
      return choose;
    } else {
      return null;
    }
  }

}
