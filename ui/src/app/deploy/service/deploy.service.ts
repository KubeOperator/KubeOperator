import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';
import {Execution} from '../component/operater/execution';

@Injectable({
  providedIn: 'root'
})
export class DeployService {

  constructor() {
  }

  private executionQueue = new Subject<Execution>();
  $executionQueue = this.executionQueue.asObservable();
  private finished = new Subject<boolean>();
  $finished = this.finished.asObservable();

  next(execution: Execution) {
    this.executionQueue.next(execution);
  }

  nextState(state: boolean) {
    this.finished.next(state);
  }

}
