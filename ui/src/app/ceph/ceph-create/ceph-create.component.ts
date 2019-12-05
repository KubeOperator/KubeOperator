import {Component, OnInit} from '@angular/core';
import * as globals from '../../globals';


@Component({
  selector: 'app-ceph-create',
  templateUrl: './ceph-create.component.html',
  styleUrls: ['./ceph-create.component.css']
})
export class CephCreateComponent implements OnInit {

  opened = false;
  name_pattern = globals.host_name_pattern;
  name_pattern_tip = globals.host_name_pattern_tip;

  constructor() {
  }

  ngOnInit() {
  }

  open() {
    this.opened = true;
  }
}
