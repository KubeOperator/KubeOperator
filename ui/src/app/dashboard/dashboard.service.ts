import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {DashboardSearch} from './dashboardSearch';

@Injectable({
  providedIn: 'root'
})
export class DashboardService {

  dashboardUrl = '/api/v1/dashboard/';


  constructor(private http: HttpClient) { }

  getDashboard(project_name: string): Observable<any> {
    return this.http.get<any>(this.dashboardUrl + project_name + '/');
  }
}
