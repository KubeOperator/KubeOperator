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

  // baseUrl = '/api/v1/message/mailbox';

  constructor(http: HttpClient) {
    super(http);
  }

}
