import {Injectable, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, Subject} from 'rxjs';
import {Log} from './log';
import {WebsocketService} from './websocket.service';

const baseUrl = '/api/v1/cluster/{clusterId}/log';

@Injectable()
export class LogService implements OnInit {

  private url = 'ws:localhost:4200/ws/tasks/fe5b341d-82f9-4987-b3d9-5301ad3421d7/log/';
  messages: Subject<any>;

  constructor(private http: HttpClient, private wsService: WebsocketService) {
    this.messages = this.wsService.connect(this.url);
  }

  getLogs(clusterId): Observable<Log[]> {
    return this.http.get<Log[]>(`${baseUrl.replace('{clusterId}', clusterId)}`);
  }

  ngOnInit(): void {

  }


}
