import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Item} from './item';
import {ItemService} from './item.service';
import {Observable} from 'rxjs';
import {Cluster} from '../cluster/cluster';
import {map, take} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class ItemRoutingResolverService implements Resolve<Item> {

  constructor(private itemService: ItemService) {
  }

  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Item> {
    const itemName = route.params['itemName'];
    return this.itemService.getItem(itemName).pipe(
      take(1),
      map(item => {
        if (item) {
          return item;
        } else {
          return null;
        }
      })
    );
  }
}
