import {Component, OnInit} from '@angular/core';
import {OfflineService} from '../../offline/offline.service';
import {Offline} from '../../offline/Offline';

@Component({
  selector: 'app-cluster-config',
  templateUrl: './cluster-config.component.html',
  styleUrls: ['./cluster-config.component.css']
})
export class ClusterConfigComponent implements OnInit {
  offlines: Offline[] = [];

  constructor(private offlineService: OfflineService) {
  }

  ngOnInit() {
    this.listOffline();
  }

  listOffline() {
    this.offlineService.listOfflines().subscribe(data => this.offlines = data);
  }

}
