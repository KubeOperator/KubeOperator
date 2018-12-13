import { Injectable } from '@angular/core';
import { Router, Resolve, RouterStateSnapshot, ActivatedRouteSnapshot } from '@angular/router';
import { Project } from './project';
import { ProjectService } from "./project.service";
import { Observable } from "rxjs/index";
import { tap } from "rxjs/internal/operators";


@Injectable()
export class ProjectRoutingResolver implements Resolve<Project> {

  constructor(private service: ProjectService, private router: Router) { }

  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Project> {
    const projectName = route.params.project;
    if (!projectName) {
      this.router.navigate(['/ansible', 'projects'])
    }
    return this.service.getProject(projectName).pipe(
      tap(project => {
        if (!project) {
          this.router.navigate(['/ansible', 'projects'])
        }
      })
    );
  }
}
