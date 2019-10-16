import {Component, OnInit} from '@angular/core';
import {Nfs} from '../nfs';

@Component({
  selector: 'app-nfs-list',
  templateUrl: './nfs-list.component.html',
  styleUrls: ['./nfs-list.component.css']
})
export class NfsListComponent implements OnInit {

  constructor() {
  }

  loading = false;
  selected: Nfs[] = [];
  items: Nfs[] = [];

  ngOnInit() {
  }

}
