import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';

import { BaseService } from '../base/base.service';
import { Playbook, PlaybookExecution } from './playbook';

@Injectable({
  providedIn: 'root'
})
export class PlaybookService extends BaseService {
  projectListUrl = '/api/v1/projects';

  constructor(private http: HttpClient) {
    super();
  }

  getPlaybooks(projectName): Observable<Playbook[]> {
    const url = `${this.projectListUrl}/${projectName}/playbooks`;
    return this.http.get<Playbook[]>(url);
  }

  getPlaybookDetail(playbook): Observable<Playbook> {
    const url = `${this.projectListUrl}/${playbook.project}/playbooks/${playbook.id}`;
    return this.http.get<Playbook>(url);
  }

  executePlaybook(playbook: Playbook): Observable<PlaybookExecution> {
    const url = `${this.projectListUrl}/${playbook.project}/playbooks/executions/`;
    return this.http.post<PlaybookExecution>(url, {'playbook': playbook.id}).pipe(
      tap(() => this.log(`Execute playbook name=${playbook.name}`)),
      catchError(this.handleError<any>('executePlaybook'))
    );
  }

  createPlaybook(playbook: Playbook, projectName: string): Observable<Playbook> {
    const url = `${this.projectListUrl}/${projectName}/playbooks/`;
    return this.http.post<Playbook>(url, playbook).pipe(
      tap(() => this.log(`Create playbook name=${playbook.name}`),),
      catchError(this.handleError<any>('executePlaybook'), )
    );
  }
}
