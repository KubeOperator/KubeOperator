import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from "@angular/router";
import {MultiClusterRepository} from "./multi-cluster-repository";
import {ProjectService} from "../project/project.service";
import {Observable} from "rxjs";
import {Project} from "../project/project";
import {map, take} from "rxjs/operators";
import {MultiClusterRepositoryService} from "./multi-cluster-repository.service";

@Injectable({
    providedIn: 'root'
})
export class MultiClusterRoutingResolverService implements Resolve<MultiClusterRepository> {
    constructor(private service: MultiClusterRepositoryService) {
    }

    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<MultiClusterRepository> | Promise<MultiClusterRepository> | MultiClusterRepository {
        const repoName = route.params.name;
        return this.service.get(repoName).pipe(
            take(1),
            map(repo => {
                if (repo) {
                    return repo;
                } else {
                    return null;
                }
            })
        );
    }
}
