import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Observable} from 'rxjs';
import {Project} from '../../project/project';
import {map, take} from 'rxjs/operators';
import {IpPoolService} from './ip-pool.service';
import {IpPool} from './ip-pool';

@Injectable({
    providedIn: 'root'
})
export class IpPoolRoutingResolverService implements Resolve<IpPool> {

    constructor(private service: IpPoolService) {
    }

    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<IpPool> | Promise<IpPool> | IpPool {
        const poolName = route.params.name;
        return this.service.get(poolName).pipe(
            take(1),
            map(ipPool => {
                if (ipPool) {
                    return ipPool;
                } else {
                    return null;
                }
            })
        );
    }
}
