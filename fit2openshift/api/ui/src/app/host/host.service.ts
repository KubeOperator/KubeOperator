import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BaseService } from '../base/base.service';

import { Host } from './host';
import { Observable } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class HostService extends BaseService {
  projectListUrl = '/api/v1/projects';

  constructor(private http: HttpClient) {
    super();
  }

  getHosts(projectName): Observable<Host[]> {
    const url = `${this.projectListUrl}/${projectName}/inventory/hosts/`;
    return this.http.get<Host[]>(url);
  }

  createHost(projectName: string, host: Host): Observable<Host> {
    const url = `${this.projectListUrl}/${projectName}/inventory/hosts/`;
    return this.http.post<Host>(url, host).pipe(
      tap(() => this.log(`Create host name=${host.name}`)),
      catchError(this.handleError<any>('createHost'))
    );
  }
}
