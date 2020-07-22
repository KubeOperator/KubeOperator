import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Project} from './project';
import {Observable} from 'rxjs';
import {ProjectService} from './project.service';
import {map, take} from 'rxjs/operators';

@Injectable({
    providedIn: 'root'
})
export class ProjectRoutingResolverService implements Resolve<Project> {

    constructor(private service: ProjectService) {
    }

    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Project> | Promise<Project> | Project {
        const projectName = route.params.name;
        return this.service.get(projectName).pipe(
            take(1),
            map(project => {
                if (project) {
                    return project;
                } else {
                    return null;
                }
            })
        );
    }
}
