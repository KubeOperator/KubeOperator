import {Injectable, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Execution} from '../deploy/component/operater/execution';


const baseUrl = '/api/v1/clusters/{clusterName}/executions/';

export class TaskLog {
  data: string;
  end: false;
}

@Injectable()
export class LogService implements OnInit {


  constructor(private http: HttpClient) {

  }

  ngOnInit(): void {

  }

  listExecutions(clusterName): Observable<Execution[]> {
    return this.http.get<Execution[]>(`${baseUrl.replace('{clusterName}', clusterName)}`);
  }

  getExecutionLog(url): Observable<TaskLog> {
    return this.http.get<TaskLog>(url);
  }

}
