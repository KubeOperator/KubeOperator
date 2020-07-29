import {Injectable} from '@angular/core';
import {ProjectService} from '../project/project.service';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Observable} from 'rxjs';
import {Project} from '../project/project';
import {map, take} from 'rxjs/operators';

@Injectable({
    providedIn: 'root'
})
export class ProjectClusterResolverService implements Resolve<Project> {

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
