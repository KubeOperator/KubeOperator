import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { Observable, of } from "rxjs";
import { catchError, map, tap } from 'rxjs/operators';
import { Role } from './role';

@Injectable({
  providedIn: 'root'
})
export class RoleService {
  projectListUrl = '/api/v1/projects';

  constructor(private http: HttpClient) { }

  getRoles(projectName): Observable<Role[]> {
    const url = `${this.projectListUrl}/${projectName}/roles`;
    return this.http.get<Role[]>(url)
      .pipe(
        tap(() => this.log('Fetch projects')),
        catchError(this.handleError('getRoles', []))
      )
  }

  getRole(projectName, roleId): Observable<Role> {
    const url = `${this.projectListUrl}/${projectName}/roles/${roleId}`;
    return this.http.get<Role>(url)
      .pipe(
        tap(() => this.log('Fetch project detail')),
        catchError(this.handleError('getRoleDetail', null))
      )
  }

  private log(msg) {
    console.log(msg)
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      this.log(`${operation} failed: ${error.message}`);
      return of(result as T);
    };
  }
}
