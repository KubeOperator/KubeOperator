import {Injectable} from '@angular/core';
import {Observable, Observer, Subject} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {

  constructor() {
  }

  private subject: Subject<MessageEvent>;

  connect(url): Subject<MessageEvent> {
    this.subject = this.create(url);
    return this.subject;
  }

  create(url): Subject<MessageEvent> {
    const ws = new WebSocket(url);

    const observable = Observable.create(
      (obs: Observer<MessageEvent>) => {
        ws.onmessage = obs.next.bind(obs);
        ws.onerror = obs.error.bind(obs);
        ws.onclose = obs.complete.bind(obs);
        return ws.close.bind(ws);
      });
    const observer = {
      next: (data: Object) => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify(data));
        }
      }
    };
    return Subject.create(observer, observable);
  }
}

