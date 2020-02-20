import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Cluster} from './cluster';
import {Observable} from 'rxjs';
import {ClusterService} from './cluster.service';
import {map, take, tap} from 'rxjs/operators';

@Injectable()
export class ClusterRoutingResolverService implements Resolve<Cluster> {
  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Cluster> {
    const clusterName = route.params['name'];
    const itemName = route.params['itemName'];
    return this.clusterService.getCluster(clusterName).pipe(
      take(1),
      map(cluster => {
        if (cluster) {
          cluster.item_name = itemName;
          return cluster;
        } else {
          return null;
        }
      })
    );
  }


  constructor(private clusterService: ClusterService) {
  }
}
