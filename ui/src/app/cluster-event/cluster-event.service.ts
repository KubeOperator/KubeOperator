import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {ClusterEventSearch} from './cluster-event-search';

@Injectable({
  providedIn: 'root'
})

export class ClusterEventService {

  baseUrl = '/api/v1/cluster/';

  constructor(private http: HttpClient) {
  }

  listClusterEvents(project_name: string, params: ClusterEventSearch): Observable<any> {
    return this.http.post<any>(this.baseUrl + project_name + '/event/', params);
  }
}
