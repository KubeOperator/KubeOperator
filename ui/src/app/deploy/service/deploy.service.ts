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

  next(execution: Execution) {
    this.executionQueue.next(execution);
  }
}
