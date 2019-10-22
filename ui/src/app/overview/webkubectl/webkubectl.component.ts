import {Component, OnInit} from '@angular/core';
import {DomSanitizer} from '@angular/platform-browser';

@Component({
  selector: 'app-webkubectl',
  templateUrl: './webkubectl.component.html',
  styleUrls: ['./webkubectl.component.css']
})
export class WebkubectlComponent implements OnInit {


  opened = false;
  loading = true;
  url = '';
  jump = {};

  constructor(private sanitizer: DomSanitizer) {

  }

  ngOnInit() {
  }

  open() {
    this.opened = true;
    this.jump = this.sanitizer.bypassSecurityTrustResourceUrl(this.url); // 信任该url
    this.loading = false;
  }

  close() {
    this.opened = false;
  }

}
