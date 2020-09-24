import { Component, OnInit } from '@angular/core';
import { Notice } from '../notice';

@Component({
  selector: 'app-mailbox-detail',
  templateUrl: './mailbox-detail.component.html',
  styleUrls: ['./mailbox-detail.component.css']
})
export class MailboxDetailComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
  }

  open(item: Notice) {
    //
  }

}
