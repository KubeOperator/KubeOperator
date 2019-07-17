import {Injectable} from '@angular/core';
import {ClusterService} from '../cluster/cluster.service';
import {Cluster} from '../cluster/cluster';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class OverviewService {

  constructor(private clusterService: ClusterService, private http: HttpClient) {

  }

  downLoad(cluster: Cluster) {
    window.open('/api/v1/cluster/' + cluster.id + '/download/');
  }
}
