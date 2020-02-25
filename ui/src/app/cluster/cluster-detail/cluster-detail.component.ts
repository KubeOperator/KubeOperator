import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Cluster} from '../cluster';
import {SessionService} from '../../shared/session.service';

@Component({
  selector: 'app-cluster-detail',
  templateUrl: './cluster-detail.component.html',
  styleUrls: ['./cluster-detail.component.css']
})
export class ClusterDetailComponent implements OnInit {

  currentCluster: Cluster;
  permission;

  constructor(private router: Router, private route: ActivatedRoute, private sessionService: SessionService) {
  }

  ngOnInit() {
    this.route.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.permission = this.sessionService.getItemPermission(this.currentCluster.item_name);
    });
  }

  backToCluster() {
    this.router.navigate(['cluster']);
  }

}
