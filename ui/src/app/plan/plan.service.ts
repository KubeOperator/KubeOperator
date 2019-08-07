import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Plan} from './plan';

@Injectable({
  providedIn: 'root'
})
export class PlanService {
  baseUrl = '/api/v1/plans/';

  constructor(private http: HttpClient) {
  }

  listPlan(): Observable<Plan[]> {
    return this.http.get<Plan[]>(this.baseUrl);
  }

  createPlan(item: Plan): Observable<Plan> {
    return this.http.post<Plan>(this.baseUrl, item);
  }

  deletePlan(name: string): Observable<Plan> {
    return this.http.delete<Plan>(this.baseUrl + name + '/');
  }

}
