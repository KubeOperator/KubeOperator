import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Router} from '@angular/router';

@Component({
  selector: 'app-local-mail-detail',
  templateUrl: './local-mail-detail.component.html',
  styleUrls: ['./local-mail-detail.component.css']
})
export class LocalMailDetailComponent implements OnInit {

  open = false;
  message: any;

  constructor(private router: Router) {
  }

  ngOnInit() {
  }

  cancel() {
    this.open = false;
  }

  toPage(url) {
    this.redirect(url);
  }


  redirect(url: string) {
    if (url) {
      const linkUrl = [url];
      this.router.navigate(linkUrl);
    }
  }
}
