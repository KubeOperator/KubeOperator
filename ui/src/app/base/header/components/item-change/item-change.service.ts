import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';
import {Profile} from '../../../../shared/session-user';

@Injectable({
  providedIn: 'root'
})
export class ItemChangeService {

  subject = new Subject<Profile>();
  public $noticeChannel = this.subject.asObservable();

  constructor() {
  }

}
