import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-cluster-event-detail',
  templateUrl: './cluster-event-detail.component.html',
  styleUrls: ['./cluster-event-detail.component.css']
})
export class ClusterEventDetailComponent implements OnInit {

  open = false;
  event: any;

  constructor() {
  }

  ngOnInit() {
  }

  cancel() {
    this.open = false;
  }
}
