import { Injectable } from '@angular/core';
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
  private notifications = 0;
  private alerts = 0;

  addItem(item: Notice) {
    this.messages.push(item);
    if (item.level === 'alert') {
      this.alerts += 1;
    } else if (item.level === 'info') {
      this.notifications += 1;
    }
  }

  deleteItem(item: Notice) {
    const index = this.messages.indexOf(item);
    if (index > -1) {
      if (item.level === 'alert') {
        this.alerts -= 1;
      } else if (item.level === 'info') {
        this.notifications -= 1;
      }
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

  getNotifications() {
    return this.notifications;
  }

  getAlerts() {
    return this.alerts;
  }

  // baseUrl = '/api/v1/message/mailbox';

  constructor(http: HttpClient) {
    super(http);
  }

}
