import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../cluster/cluster';

@Component({
  selector: 'app-cluster-event',
  templateUrl: './cluster-event.component.html',
  styleUrls: ['./cluster-event.component.css']
})
export class ClusterEventComponent implements OnInit {
  currentCluster: Cluster;

  constructor(private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
    });
  }

}
