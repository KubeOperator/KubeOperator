import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ClusterEventSearch} from '../cluster-event-search';
import {ClusterEventService} from '../cluster-event.service';
import {ClusterEventDetailComponent} from '../cluster-event-detail/cluster-event-detail.component';

@Component({
  selector: 'app-cluster-event-list',
  templateUrl: './cluster-event-list.component.html',
  styleUrls: ['./cluster-event-list.component.css']
})
export class ClusterEventListComponent implements OnInit {
  @Input() currentCluster: Cluster;
  search = new ClusterEventSearch();
  items = [];
  totalItems: number;
  loading = true;
  @ViewChild(ClusterEventDetailComponent, {static: true})
  detail: ClusterEventDetailComponent;

  constructor(private clusterEventService: ClusterEventService) {
  }

  ngOnInit() {
    this.search.limitDays = 7;
    this.search.currentPage = 1;
    this.search.size = 10;
    this.search.type = 'Warning';
    this.listClusterEvent(10);
  }

  listClusterEvent(pageSize) {
    this.loading = true;
    this.search.size = pageSize;
    this.clusterEventService.listClusterEvents(this.currentCluster.name, this.search).subscribe(res => {
      this.items = res.items;
      this.totalItems = res.total;
      this.loading = false;
    });
  }

  showDetail(item: any) {
    this.detail.event = item;
    this.detail.open = true;
  }
}
