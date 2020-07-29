import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Cluster} from './cluster';
import {Observable} from 'rxjs';
import {ClusterService} from './cluster.service';
import {map, take} from 'rxjs/operators';

@Injectable({
    providedIn: 'root'
})
export class ClusterRoutingResolverService implements Resolve<Cluster> {

    constructor(private service: ClusterService) {
    }

    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Cluster> | Promise<Cluster> | Cluster {
        const clusterName = route.params.name;
        const projectName = route.params.projectName;
        return this.service.get(clusterName).pipe(
            take(1),
            map(cluster => {
                if (cluster) {
                    cluster.projectName = projectName;
                    return cluster;
                } else {
                    return null;
                }
            })
        );
    }
}
