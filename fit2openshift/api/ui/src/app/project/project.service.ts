import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, of, throwError as observableThrowError } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';

import { Project } from './project';

const httpOptions = {
  headers: new HttpHeaders({ 'Content-Type': 'application/json' })
};

@Injectable({
  providedIn: 'root'
})
export class ProjectService {
  private projectListUrl = `/api/v1/projects`;

  constructor( private http: HttpClient) { }

  getProjects(): Observable<Project[]> {
    return this.http.get<Project[]>(this.projectListUrl)
      .pipe(
        tap(() => this.log('Fetch projects')),
        catchError(this.handleError('getProjects', []))
      );
  }

  getProject(projectName): Observable<Project> {
    const projectDetailUrl = `${this.projectListUrl}/${projectName}`;
    return this.http.get<Project>(projectDetailUrl)
      .pipe(
        tap(() => this.log('Fetch project detail'))
      );
  }

  createProject(project): Observable<Project> {
    const projectCreateUrl = `${this.projectListUrl}/`;
    return this.http.post<Project>(projectCreateUrl, project).pipe(
      tap(() => this.log(`Create project id=${project.name}`)),
      catchError(this.handleError<any>('updateHero'))
    );
  }

  checkProjectExists(projectName): Observable<boolean> {
    const projectFilterUrl = `${this.projectListUrl}/?name=${projectName}`;
    return this.http.get<Project[]>(projectFilterUrl)
      .pipe(
        map(response => response.length > 0),
        catchError(error => observableThrowError(error)), );
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      this.log(`${operation} failed: ${error.message}`);
      return of(result as T);
    };
  }

  private log(msg) {
    console.log(msg);
  }

}
