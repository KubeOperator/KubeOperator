import { EventEmitter, Injectable, Output } from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {HttpClient} from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { BaseModelService } from '../../../shared/class/BaseModelService';
import { Notice } from './notice';

@Injectable({
  providedIn: 'root'
})
export class NoticeService extends BaseModelService<any> {

  private messages: Notice[] = [];

  private unreadSource = new BehaviorSubject({unreadInfo: 0, unreadAlert: 0});
  currentUnread = this.unreadSource.asObservable();

  // can be done through http
  addItem(item: Notice) {
    this.messages.push(item);
  }

  deleteItem(item: Notice) {
    const index = this.messages.indexOf(item);
    if (index > -1) {
      this.messages.splice(index, 1);
    }
    else {
      alert('oops, something goes wrong');
    }
  }

  deleteItems(items: Notice[]) {
    for (const item of items) {
      this.deleteItem(item);
    }
  }

  getItems() {
    return this.messages;
  }

  updateItemOnRead(item: Notice) {
    const index = this.messages.indexOf(item);
    if (index > -1) {
      this.messages[index].isRead = true;
    } else { console.log('can\'t update item\'s isRead property'); }
  }

  updateItemsOnRead(items: Notice[]) {
    for (const item of items) {
      this.updateItemOnRead(item);
    }
  }
  //

  // for updating the header component with the latest unread notices
  changeUnread(unread: {unreadInfo: number, unreadAlert: number}) {
    this.unreadSource.next(unread);
  }




  // baseUrl = '/api/v1/message/mailbox';

  constructor(http: HttpClient) {
    super(http);
  }

}
