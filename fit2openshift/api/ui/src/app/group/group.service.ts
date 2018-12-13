import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { catchError, tap } from 'rxjs/operators';

import { BaseService } from '../base/base.service';
import { Group } from './group';

@Injectable({
  providedIn: 'root'
})
export class GroupService extends BaseService {
  projectListUrl = '/api/v1/projects';

  constructor(private http: HttpClient) {
    super();
  }

  getGroups(projectName: string): Observable<Group[]> {
    const url = `${this.projectListUrl}/${projectName}/inventory/groups/`;
    return this.http.get<Group[]>(url).pipe(
      tap(() => this.log(`Get ${projectName} groups`)),
      catchError(this.handleError<any>('Get groups'))
    );
  }

  createGroup(projectName: string, group: Group): Observable<Group> {
    const url = `${this.projectListUrl}/${projectName}/inventory/groups/`;
    return this.http.post<Group>(url, group).pipe(
      tap(() => this.log(`Create ${projectName} groups`)),
      catchError(this.handleError<any>('Create groups'))
    );
  }
}
