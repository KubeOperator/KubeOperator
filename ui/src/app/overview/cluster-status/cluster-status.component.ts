import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {NodeService} from '../../node/node.service';
import {Node} from '../../node/node';
import {ScaleComponent} from '../scale/scale.component';
import {Router} from '@angular/router';
import {OperaterService} from '../../deploy/component/operater/operater.service';
import {ClusterHealthService} from '../../cluster-health/cluster-health.service';
import {AddWorkerComponent} from '../add-worker/add-worker.component';
import {RemoveWorkerComponent} from '../remove-worker/remove-worker.component';
import {ClusterService} from '../../cluster/cluster.service';
import {SessionService} from '../../shared/session.service';
import {stringify} from '@angular/compiler/src/util';


@Component({
    selector: 'app-cluster-status',
    templateUrl: './cluster-status.component.html',
    styleUrls: ['./cluster-status.component.css']
})
export class ClusterStatusComponent implements OnInit {

    @Input() currentCluster: Cluster;
    workers: Node[] = [];
    @ViewChild(ScaleComponent, {static: true}) scale: ScaleComponent;
    @ViewChild(AddWorkerComponent, {static: true}) addWorker: AddWorkerComponent;
    @ViewChild(RemoveWorkerComponent, {static: true}) removeWorker: RemoveWorkerComponent;
    componentData = [];
    loading = false;
    permission;

    constructor(private nodeService: NodeService, private clusterHealthService: ClusterHealthService,
                private router: Router, private operaterService: OperaterService, private clusterService: ClusterService,
                private sessionService: SessionService) {
    }

    ngOnInit() {
        this.permission = this.sessionService.getItemPermission(this.currentCluster.item_name);
        this.nodeService.listNodes(this.currentCluster.name).subscribe(data => {
            this.workers = data.filter((node) => {
                return node.roles.includes('worker');
            });
        });
        this.getClusterStatus();
    }

    handleScale() {
        const params = {'num': this.scale.worker_size};
        this.operaterService.executeOperate(this.currentCluster.name, 'scale', params).subscribe(() => {
            this.redirect('deploy');
        }, error => {
            this.scale.opened = false;
        });
    }

    refresh() {
        this.clusterService.getCluster(this.currentCluster.name).subscribe(data => {
            this.currentCluster = data;
        });
    }

    handleAddWorker() {
        const hosts = [];
        this.addWorker.host_names.forEach(h => {
            hosts.push(h['value']);
        });
        const params = {'hosts': hosts};
        this.operaterService.executeOperate(this.currentCluster.name, 'add-worker', params).subscribe(() => {
            this.redirect('deploy');
        }, error => {
            this.scale.opened = false;
        });
    }

    handleRemoveWorker() {
        const nodes = [];
        this.removeWorker.worker_names.forEach(w => {
            nodes.push(w['value']);
        });
        const params = {'nodes': nodes};
        this.operaterService.executeOperate(this.currentCluster.name, 'remove-worker', params).subscribe(() => {
            this.redirect('deploy');
        }, error => {
            this.scale.opened = false;
        });
    }

    redirect(url: string) {
        if (url) {
            const linkUrl = ['cluster', this.currentCluster.name, url];
            this.router.navigate(linkUrl);
        }
    }

    onScale() {
        this.scale.worker_size = this.workers.length;
        this.scale.opened = true;
    }

    onAddWorker() {
        this.addWorker.loadHosts();
        this.addWorker.opened = true;
    }

    onRemoveWorker() {
        this.removeWorker.loadNodes(this.currentCluster.name);
        this.removeWorker.opened = true;
    }

    toHealth() {
        this.redirect('health');
    }

    getClusterStatus() {
        this.loading = true;
        this.clusterHealthService.listComponent(this.currentCluster.name).subscribe(res => {
            this.componentData = res;
            this.loading = false;
        }, error1 => {
            this.loading = false;
        });
    }


}
