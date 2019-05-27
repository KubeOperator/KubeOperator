import {Injectable} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class DeplayUtilService {

  constructor() {
  }

  private compated_status = ['SUCCESS', 'FAILURE'];

  execution_is_complated(status: string): boolean {
    return this.compated_status.includes(status);
  }

}
