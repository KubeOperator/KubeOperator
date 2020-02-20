import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Cluster} from '../cluster';

@Component({
  selector: 'app-cluster-detail',
  templateUrl: './cluster-detail.component.html',
  styleUrls: ['./cluster-detail.component.css']
})
export class ClusterDetailComponent implements OnInit {

  currentCluster: Cluster;

  constructor(private router: Router, private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.data.subscribe(data => {
      this.currentCluster = data['cluster'];
    });
  }

  backToCluster() {
    this.router.navigate(['item/' + this.currentCluster.item_name + '/cluster'], {queryParams: {name: this.currentCluster.item_name}});
  }

}
