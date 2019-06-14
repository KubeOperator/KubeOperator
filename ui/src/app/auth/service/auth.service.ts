import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Router} from '@angular/router';
import {Observable} from 'rxjs';
import {AuthTemplate} from '../class/auth';
import {Cluster, ExtraConfig} from '../../cluster/cluster';
import {ClusterService} from '../../cluster/cluster.service';
import {OperaterService} from '../../deploy/component/operater/operater.service';

export const url = '/api/v1/auth/';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  constructor(private http: HttpClient, private clusterService: ClusterService,
              private optService: OperaterService, private router: Router) {

  }

  listAuthTemplate(): Observable<AuthTemplate[]> {
    return this.http.get<AuthTemplate[]>(url);
  }

  getAuthTemplate(name: string): Observable<AuthTemplate> {
    return this.http.get<AuthTemplate>(url + name);
  }

  fullAuth(auth: AuthTemplate, clusterName: string) {
    this.clusterService.getClusterConfig(clusterName, 'openshift_master_identity_providers').subscribe(data => {
      const d = data.value[0];
      auth.meta.options.forEach(option => {
        for (const key in d) {
          if (key === option.name) {
            option.value = d[key];
          }
        }
      });
      auth.meta.vars.forEach(v => {
        this.clusterService.getClusterConfig(clusterName, v.name).subscribe(vs => {
          v.value = vs.value;
        });
      });
      console.log(auth);
    });
  }

  configAuth(auth: AuthTemplate, cluster: Cluster) {
    const config = auth.meta.config;
    auth.meta.options.forEach(option => {
      config[option.name] = option.value;
    });
    const auth_list = [];
    auth_list.push(config);
    const authConfig: ExtraConfig = new ExtraConfig();
    authConfig.key = 'openshift_master_identity_providers';
    authConfig.value = auth_list;
    const promises: Promise<{}>[] = [];
    const vars: ExtraConfig[] = [];
    auth.meta.vars.forEach(_var => {
      const c: ExtraConfig = new ExtraConfig();
      c.key = _var.name;
      c.value = _var.value;
      vars.push(c);
    });
    console.log(authConfig);
    promises.push(this.clusterService.configCluster(cluster.name, authConfig).toPromise());
    vars.forEach(c => {
      promises.push(this.clusterService.configCluster(cluster.name, c).toPromise());
    });
    Promise.all(promises).then(nil => {
      this.clusterService.configClusterAuth(cluster.name, auth.name).subscribe(d => {
        console.log(auth.name);
        // this.optService.executeOperate(cluster.name, 'config-auth').subscribe(data => {
        //   this.router.navigate(['kubeOperator', 'cluster', cluster.name, 'deploy']);
        // });
      });
    });
  }
}
