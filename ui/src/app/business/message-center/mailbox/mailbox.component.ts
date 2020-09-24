import { Component, OnInit, ViewChild } from '@angular/core';
import { MailboxListComponent } from './mailbox-list/mailbox-list.component';
import { MailboxDetailComponent } from './mailbox-detail/mailbox-detail.component';
import { MailboxCreateComponent } from './mailbox-create/mailbox-create.component';
import { MailboxDeleteComponent } from './mailbox-delete/mailbox-delete.component';
import { Notice } from './notice';


@Component({
  selector: 'app-mailbox',
  templateUrl: './mailbox.component.html',
  styleUrls: ['./mailbox.component.css']
})
export class MailboxComponent implements OnInit {

  @ViewChild(MailboxListComponent, {static: true})
  list: MailboxListComponent;

  @ViewChild(MailboxCreateComponent, {static: true})
  create: MailboxCreateComponent;

  @ViewChild(MailboxDeleteComponent, {static: true})
  delete: MailboxDeleteComponent;

  @ViewChild(MailboxDetailComponent, {static: true})
  detail: MailboxDetailComponent;

  constructor() {
  }

  ngOnInit(): void {
  }

  openCreate() {
    this.create.open();
  }

  openDelete(items: Notice[]) {
    this.delete.open(items);
  }

  openDetail(item: Notice) {
    this.detail.open(item);
  }

  refresh() {
    this.list.reset();
    this.list.refresh();
  }

}
