import {Component, OnInit} from '@angular/core';
import {Cluster} from '../cluster/cluster';
import {ActivatedRoute} from '@angular/router';

@Component({
  selector: 'app-config',
  templateUrl: './config.component.html',
  styleUrls: ['./config.component.css']
})
export class ConfigComponent implements OnInit {

  currentCluster: Cluster;

  constructor(private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
    });
  }

}
