import { Injectable } from '@angular/core';
import {Observable} from 'rxjs';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ClusterStorageService {
  baseUrl = '/api/v1/cluster/';

  constructor(private http: HttpClient) { }

  listClusterStorage(project_name: string): Observable<any> {
    return this.http.get<any>(this.baseUrl + project_name + '/storage/');
  }
}
