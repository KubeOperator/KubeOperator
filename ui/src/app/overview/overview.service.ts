import {Injectable} from '@angular/core';
import {ClusterService} from '../cluster/cluster.service';
import {Cluster} from '../cluster/cluster';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {ClusterToken} from './describe/class/describe';

@Injectable({
  providedIn: 'root'
})
export class OverviewService {

  constructor(private clusterService: ClusterService, private http: HttpClient) {

  }

  downLoad(cluster: Cluster) {
    window.open('/api/v1/cluster/' + cluster.id + '/download/');
  }

  getClusterToken(cluster: Cluster): Observable<ClusterToken> {
    return this.http.get<ClusterToken>('/api/v1/cluster/' + cluster.id + '/token/');
  }
}
