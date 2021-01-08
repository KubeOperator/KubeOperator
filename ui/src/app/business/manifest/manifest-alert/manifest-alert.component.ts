import { Component, OnInit } from '@angular/core';
import {Manifest} from '../manifest';

@Component({
  selector: 'app-manifest-alert',
  templateUrl: './manifest-alert.component.html',
  styleUrls: ['./manifest-alert.component.css']
})
export class ManifestAlertComponent implements OnInit {

  constructor() {
  }

  opened = false;
  ngOnInit(): void {
  }

  close() {
    this.opened = false;
  }

  open() {
    this.opened = true;
  }
}
