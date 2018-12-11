import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Cluster} from './cluster';
import {Observable} from 'rxjs';

@Injectable()
export class ClusterRoutingResolverService implements Resolve<Cluster> {
  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Cluster> | Promise<Cluster> | Cluster {
    const clusterId = route.params['id'] == null ? route.queryParams['clusterId'] : route.params['id'];
    return undefined;
  }


  constructor() {
  }
}
