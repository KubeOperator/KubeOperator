import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-local-mail-detail',
  templateUrl: './local-mail-detail.component.html',
  styleUrls: ['./local-mail-detail.component.css']
})
export class LocalMailDetailComponent implements OnInit {

  open = false;
  message: any;

  constructor() {
  }

  ngOnInit() {
    console.log(this.message);
  }

  cancel() {
    this.open = false;
  }
}
