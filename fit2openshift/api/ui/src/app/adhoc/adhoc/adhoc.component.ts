import { Component, OnInit } from '@angular/core';

import { Adhoc } from './adhoc';

@Component({
  selector: 'app-adhoc',
  templateUrl: './adhoc.component.html',
  styles: []
})
export class AdhocComponent implements OnInit {
  adhoc: Adhoc;

  constructor() { }

  ngOnInit() {
    this.adhoc = new Adhoc();
    this.adhoc.module = 'shell';
  }

}
