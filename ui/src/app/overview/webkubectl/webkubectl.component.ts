import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-webkubectl',
  templateUrl: './webkubectl.component.html',
  styleUrls: ['./webkubectl.component.css']
})
export class WebkubectlComponent implements OnInit {

  constructor() {
  }

  opened = false;
  loading = true;
  url = '';

  ngOnInit() {
  }

}
